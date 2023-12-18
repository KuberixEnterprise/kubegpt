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

// SerializeObjectAsJSON serializes an object as JSON and returns the result
func SerializeObjectAsJSON(ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object, eventResource v1alpha1.Event) (v1alpha1.Result, v1alpha1.Store, error) {
	var result v1alpha1.Result
	var store v1alpha1.Store

	if err := c.Get(ctx, key, obj); err != nil {
		return result, store, err
	}

	if pod, ok := obj.(*v1.Pod); ok {
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
		return result, store, nil
	}

	if service, ok := obj.(*v1.Service); ok {
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
