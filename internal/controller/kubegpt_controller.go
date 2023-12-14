/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"time"

	c "github.com/kuberixenterprise/kubegpt/pkg/cache"

	"github.com/kuberixenterprise/kubegpt/pkg/ai"
	"github.com/kuberixenterprise/kubegpt/pkg/integrations"
	"github.com/kuberixenterprise/kubegpt/pkg/resource"
	"github.com/kuberixenterprise/kubegpt/pkg/sinks"
	"github.com/sirupsen/logrus"
	v1app "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/events/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1alpha1 "github.com/kuberixenterprise/kubegpt/api/v1alpha1"
)

// KubegptReconciler reconciles a Kubegpt object
type KubegptReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	Integrations *integrations.Integrations
}

//+kubebuilder:rbac:groups=core.kubegpt.io,resources=kubegpts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.kubegpt.io,resources=kubegpts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.kubegpt.io,resources=kubegpts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kubegpt object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *KubegptReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	l := log.FromContext(ctx)
	// result list 선언 + 에러 확인
	kubegptConfig := &corev1alpha1.Kubegpt{}
	if err := r.Client.Get(ctx, req.NamespacedName, kubegptConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	cacheFilePath := kubegptConfig.Spec.Cache.Path

	events := &v12.EventList{}
	if err := r.Client.List(ctx, events); err != nil {
		return ctrl.Result{}, err
	}

	resultList := &corev1alpha1.ResultList{}

	for _, event := range events.Items {
		if event.Type == "Warning" {
			eventResource := corev1alpha1.Event{
				Type:    event.Type,
				Reason:  event.Reason,
				Count:   int16(event.DeprecatedCount),
				Message: event.Note,
			}
			var err error
			var obj client.Object
			switch event.Regarding.Kind {
			case "Pod":
				pod := &v1.Pod{}
				err = r.Get(ctx, client.ObjectKey{Name: event.Regarding.Name, Namespace: event.Regarding.Namespace}, pod)
				obj = pod
			case "Service":
				service := &v1.Service{}
				err = r.Get(ctx, client.ObjectKey{Name: event.Regarding.Name, Namespace: event.Regarding.Namespace}, service)
				obj = service
			case "Deployment":
				deployment := &v1app.Deployment{}
				err = r.Get(ctx, client.ObjectKey{Name: event.Regarding.Name, Namespace: event.Regarding.Namespace}, deployment)
			case "Replicaset":
				replicaset := &v1app.ReplicaSet{}
				err = r.Get(ctx, client.ObjectKey{Name: event.Regarding.Name, Namespace: event.Regarding.Namespace}, replicaset)
			default:
				continue
			}
			if err != nil {
				logrus.Info("조회 실패 이벤트를 제외합니다.", "name", event.Regarding.Name, "namespace", event.Regarding.Namespace)
				continue
			} else {
				jsonString, store, err := resource.SerializeObjectAsJSON(ctx, r.Client, client.ObjectKey{Name: event.Regarding.Name, Namespace: event.Regarding.Namespace}, obj, eventResource)
				if err != nil {
					// 에러 처리
					continue
				}
				resultList.Items = append(resultList.Items, jsonString)
				resultList.Store = append(resultList.Store, store)
			}

		}
	}

	if kubegptConfig.Spec.Sink != nil && kubegptConfig.Spec.Sink.Type != "" && kubegptConfig.Spec.Sink.Endpoint != "" {
		// sink 설정
		var slackSink sinks.SlackSink
		slackSink.Configure(kubegptConfig)
		cache := c.NewCache()
		err := cache.LoadCacheFromFile(cacheFilePath)
		if err != nil {
			logrus.Error(err, "캐시 읽기 실패")
			return ctrl.Result{}, err
		}
		logrus.Info("캐시 파일을 읽어 옵니다.", "Load Cache Count ", len(cache.Data))
		var keystore []string
		for _, result := range resultList.Items {
			var res corev1alpha1.Result
			result := result
			key := result.Spec.Name + "_" + result.Spec.Namespace + "_" + result.Spec.Kind + "_" + result.Spec.Event[0].Message
			value := result.Spec.Event[0].Reason
			count := int(result.Spec.Event[0].Count)
			keystore = append(keystore, key)
			if !cache.DuplicateEvent(key, value) {
				cache.CacheAdd(key, value, count)
				cache.SaveCacheToFile(cacheFilePath)
				logrus.Info("캐시 데이터 저장(New Event)", "Add Cache", key)

				if err := r.Get(ctx, client.ObjectKey{Name: result.Name, Namespace: result.Namespace}, &res); err == nil {
					l.Error(err, "Result 조회 실패", "name", result.Name, "namespace", result.Namespace)
				}

				if res.Status.Webhook == "" {
					// 슬랙에 새로 보내는 로직
					gptMsg, err := slackSink.Emit(result.Spec, kubegptConfig.Spec)
					if err != nil {
						logrus.Error(err, "Sink 발송 실패")
						return ctrl.Result{}, err
					}
					result.Status.Webhook = kubegptConfig.Spec.Sink.Endpoint
					if kubegptConfig.Spec.AI.Enabled {
						go func() {
							content := gptMsg
							answer := ai.GetAnswer(content, kubegptConfig.Spec)
							answerData, err := json.Marshal(sinks.StringSlackMessage(answer, result.Spec))
							if err != nil {
								logrus.Error(err, "Failed to marshal message")
								return
							}
							sinks.SlackClient(&slackSink, answerData, "chatGPT Answer")
							cache.CacheGPTUpdate(key, answer)
							cache.SaveCacheToFile(cacheFilePath)
						}()
					}
				} else {
					res.Status.Webhook = ""
				}
			} else {
				// 캐시에 이미 있는 경우
				// 20분이 지난 경우 슬랙에 보내고 캐시 업데이트
				if time.Since(cache.Data[key].Timestamp) > time.Duration(kubegptConfig.Spec.Timer.SlackInterval)*time.Minute {
					// 20분 경과 했지만 error Count 증가 없는 경우 에러 해결로 판단 pass
					if time.Since(cache.Data[key].ErrorTime) <= time.Duration(10+kubegptConfig.Spec.Timer.SlackInterval)*time.Minute {
						// 슬랙에 새로 보내는 로직
						// 20분이 지난 경우 슬랙에 보내고 캐시 업데이트
						err := slackSink.ReEmit(key, cache.Data[key])
						if err != nil {
							logrus.Error(err, "Sink 발송 실패")
							return ctrl.Result{}, err
						}
						cache.CacheTimeUpdate(key)
						cache.SaveCacheToFile(cacheFilePath)
						logrus.Println("캐시 Timestamp 업데이트", "key", key)
					} else {
						if count > cache.Data[key].ErrorCount || count < 10 {
							// 에러 카운트가 증가한 경우 ErrorTime 업데이트
							cache.CacheErrorTimeUpdate(key, count)
						}
						// 추가 에러가 없으므로 패스
						logrus.Println("추가 에러 없으므로 패스:\n", "key", key)
					}
				} else {
					// 20분이 지나지 않은 경우
					if count > cache.Data[key].ErrorCount || count < 10 {
						// 에러 카운트가 증가한 경우 ErrorTime 업데이트
						cache.CacheErrorTimeUpdate(key, count)
					}
					logrus.Println("count", count, "cache.Data[key].ErrorCount", cache.Data[key].ErrorCount)
					logrus.Println("duplicate event pass")

				}
			}
		}
		cache.Cleanup(keystore)

	}

	// 결과 상태 업데이트
	return ctrl.Result{RequeueAfter: time.Duration(kubegptConfig.Spec.Timer.ErrorInterval) * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubegptReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Kubegpt{}).
		Watches(&v12.Event{}, &handler.EnqueueRequestForObject{}).
		Watches(&corev1alpha1.Result{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1.Pod{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1app.Deployment{}, &handler.EnqueueRequestForObject{}).
		Watches(&v1app.ReplicaSet{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
