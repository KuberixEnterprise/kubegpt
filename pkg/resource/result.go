package resource

import (
	"context"

	"github.com/kuberixenterprise/kubegpt/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// import to yaml
	//"k8s.io/apimachinery/pkg/runtime"
	////"k8s.io/apimachinery/pkg/runtime/serializer/json"
	//"bytes"
	//"k8s.io/apimachinery/pkg/types"
)

// SerializeObjectAsJSON 함수는 주어진 Kubernetes 오브젝트를 JSON 형식으로 직렬화합니다.
func SerializeObjectAsJSON(ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object, eventResource v1alpha1.Event) (v1alpha1.Result, error) {
	var result v1alpha1.Result

	// 오브젝트를 Kubernetes API 서버로부터 가져옵니다.
	if err := c.Get(ctx, key, obj); err != nil {
		return result, err
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
		//// 결과를 JSON으로 직렬화
		//jsonBytes, err := json.Marshal(result)
		//if err != nil {
		//	return "", err
		//}
		return result, nil
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
		return result, nil
	}

	return result, nil
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
