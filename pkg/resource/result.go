package resource

import (
	"context"

	"github.com/kuberixenterprise/kubegpt/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// import to yaml
	//"k8s.io/apimachinery/pkg/runtime"
	////"k8s.io/apimachinery/pkg/runtime/serializer/json"
	//"bytes"
	//"k8s.io/apimachinery/pkg/types"
)

// SerializeObjectAsJSON 함수는 주어진 Kubernetes 오브젝트를 JSON 형식으로 직렬화합니다.
func SerializeObjectAsJSON(ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object, eventResource v1alpha1.Event) (v1alpha1.Result, v1alpha1.Store, error) {
	var result v1alpha1.Result
	var store v1alpha1.Store

	// 오브젝트를 Kubernetes API 서버로부터 가져옵니다.
	if err := c.Get(ctx, key, obj); err != nil {
		return result, store, err
	}

	// Pod 타입의 리소스 처리
	if pod, ok := obj.(*v1.Pod); ok {
		// 이미지와 라벨 정보만 추출
		result = v1alpha1.Result{
			Spec: v1alpha1.ResultSpec{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				Kind:      pod.Kind,
				Images:    extractImagesFromPod(pod),
				Labels:    pod.Labels,
				Event:     []v1alpha1.Event{eventResource},
			},
		}

		store = v1alpha1.Store{
			Kind:      pod.Kind,
			Name:      pod.Name,
			Namespace: pod.Namespace,
		}
		//// 결과를 JSON으로 직렬화
		//jsonBytes, err := json.Marshal(result)
		//if err != nil {
		//	return "", err
		//}
		return result, store, nil
	}

	// 다른 타입의 리소스에 대한 처리 (필요한 경우)
	if service, ok := obj.(*v1.Service); ok {
		// 이미지와 라벨 정보만 추출
		result = v1alpha1.Result{
			Spec: v1alpha1.ResultSpec{
				Name:      service.Name,
				Namespace: service.Namespace,
				Kind:      service.Kind,
				Labels:    service.Labels,
				Event:     []v1alpha1.Event{eventResource},
			},
		}
		pods := &v1.PodList{}
		labelSelector := labels.SelectorFromSet(service.Labels)
		listOptions := &client.ListOptions{
			Namespace:     service.Namespace,
			LabelSelector: labelSelector,
		}
		if err := c.List(ctx, pods, listOptions); err != nil {
			return result, store, err
		}
		for _, pod := range pods.Items {
			store = v1alpha1.Store{
				Kind:      pod.Kind,
				Name:      pod.Name,
				Namespace: pod.Namespace,
			}
		}
		return result, store, nil
	}

	return result, store, nil
}

func extractImagesFromPod(pod *v1.Pod) []string {
	var images []string
	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}
	for _, initContainer := range pod.Spec.InitContainers {
		images = append(images, initContainer.Image)
	}
	return images
}
