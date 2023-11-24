package resource

import (
	corev1 "k8s.io/api/core/v1"
	r1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test.kubebuilder.io/project/api/v1alpha1"
)

const (
	DeploymentName = "kubegpt-deployment"
)

func GetService(config v1alpha1.Kubegpt) (*corev1.Service, error) {
	// Create service
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubegpt",
			Namespace: config.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind:       config.Kind,
					Name:       config.Name,
					UID:        config.UID,
					APIVersion: config.APIVersion,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": DeploymentName,
			},
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
				},
			},
		},
	}

	return &service, nil
}

func GetServiceAccount(config v1alpha1.Kubegpt) (*corev1.ServiceAccount, error) {
	// Create service account
	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubegpt",
			Namespace: config.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind:       config.Kind,
					Name:       config.Name,
					UID:        config.UID,
					APIVersion: config.APIVersion,
				},
			},
		},
	}

	return &serviceAccount, nil
}

func GetClusterRoleBinding(config v1alpha1.Kubegpt) (*r1.ClusterRoleBinding, error) {

	// Create cluster role binding
	clusterRoleBinding := r1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kubegpt",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: config.Kind,
					Name: config.Name,
					UID:  config.UID,
				},
			},
		},
		Subjects: []r1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "kubegpt",
				Namespace: config.Namespace,
			},
		},
		RoleRef: r1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "kubegpt",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	return &clusterRoleBinding, nil
}

//func GetDeployment(config v1alpha1.Kubegpt) (*appsv1.Deployment, error) {\
// Create deployment
//replicas := int32(1)
//deployment := appsv1.Deployment{
//	ObjectMeta: metav1.ObjectMeta{
//		Name:      DeploymentName,
//		Namespace: config.Namespace,
//		OwnerReferences: []metav1.OwnerReference{
//			{
//				Kind:       config.Kind,
//				Name:       config.Name,
//				UID:        config.UID,
//				APIVersion: config.APIVersion,
//			},
//		},
//	},
//	Spec: appsv1.DeploymentSpec{
//		Replicas: &replicas,
//		Selector: &metav1.LabelSelector{
//			MatchLabels: map[string]string{
//				"app": DeploymentName,
//			},
//		},
//		Template: corev1.PodTemplateSpec{
//			ObjectMeta: metav1.ObjectMeta{
//				Labels: map[string]string{
//					"app": DeploymentName,
//				},
//			},
//			Spec: corev1.PodSpec{
//				ServiceAccountName: "kubegpt",
//				Containers: []corev1.Container{
//					{
//						Name:            "kubegpt",
//						ImagePullPolicy: corev1.PullAlways,
//						Image:           "ghcr.io/kubegpt-ai/kubegpt:" + config.Spec.Version,
//						Args: []string{
//							"serve",
//						},
//						Env: []corev1.EnvVar{
//							{
//								Name:  "KUBEGPT_MODEL",
//								Value: config.Spec.AI.Model,
//							},
//							{
//								Name:  "KUBEGPT_BACKEND",
//								Value: config.Spec.AI.Backend,
//							},
//{
//	Name:  "XDG_CONFIG_HOME",
//	Value: "/k8sgpt-data/.config",
//},
//{
//	Name:  "XDG_CACHE_HOME",
//	Value: "/k8sgpt-data/.cache",
//},
//						},
//						Ports: []corev1.ContainerPort{
//							{
//								ContainerPort: 8080,
//							},
//						},
//						Resources: corev1.ResourceRequirements{
//							Limits: corev1.ResourceList{
//								corev1.ResourceCPU:    resource.MustParse("1"),
//								corev1.ResourceMemory: resource.MustParse("512Mi"),
//							},
//							Requests: corev1.ResourceList{
//								corev1.ResourceCPU:    resource.MustParse("0.2"),
//								corev1.ResourceMemory: resource.MustParse("156Mi"),
//							},
//						},
//						VolumeMounts: []corev1.VolumeMount{
//							{
//								MountPath: "/k8sgpt-data",
//								Name:      "k8sgpt-vol",
//							},
//						},
//					},
//				},
//				Volumes: []corev1.Volume{
//					{
//						VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
//						Name:         "k8sgpt-vol",
//					},
//				},
//			},
//		},
//	},
//}
//if config.Spec.AI.Secret != nil {
//	password := corev1.EnvVar{
//		Name: "K8SGPT_PASSWORD",
//		ValueFrom: &corev1.EnvVarSource{
//			SecretKeyRef: &corev1.SecretKeySelector{
//				LocalObjectReference: corev1.LocalObjectReference{
//					Name: config.Spec.AI.Secret.Name,
//				},
//				Key: config.Spec.AI.Secret.Key,
//			},
//		},
//	}
//	deployment.Spec.Template.Spec.Containers[0].Env = append(
//		deployment.Spec.Template.Spec.Containers[0].Env, password,
//	)
//}
//if config.Spec.RemoteCache != nil {
//
//	// check to see if key/value exists
//	addRemoteCacheEnvVar := func(name, key string) {
//		envVar := v1.EnvVar{
//			Name: name,
//			ValueFrom: &v1.EnvVarSource{
//				SecretKeyRef: &v1.SecretKeySelector{
//					LocalObjectReference: v1.LocalObjectReference{
//						Name: config.Spec.RemoteCache.Credentials.Name,
//					},
//					Key: key,
//				},
//			},
//		}
//		deployment.Spec.Template.Spec.Containers[0].Env = append(
//			deployment.Spec.Template.Spec.Containers[0].Env, envVar,
//		)
//	}
//	if config.Spec.RemoteCache.Azure != nil {
//		addRemoteCacheEnvVar("AZURE_CLIENT_ID", "azure_client_id")
//		addRemoteCacheEnvVar("AZURE_TENANT_ID", "azure_tenant_id")
//		addRemoteCacheEnvVar("AZURE_CLIENT_SECRET", "azure_client_secret")
//	} else if config.Spec.RemoteCache.S3 != nil {
//		addRemoteCacheEnvVar("AWS_ACCESS_KEY_ID", "aws_access_key_id")
//		addRemoteCacheEnvVar("AWS_SECRET_ACCESS_KEY", "aws_secret_access_key")
//	}
//}
//
//if config.Spec.AI.BaseUrl != "" {
//	baseUrl := corev1.EnvVar{
//		Name:  "K8SGPT_BASEURL",
//		Value: config.Spec.AI.BaseUrl,
//	}
//	deployment.Spec.Template.Spec.Containers[0].Env = append(
//		deployment.Spec.Template.Spec.Containers[0].Env, baseUrl,
//	)
//}
// Engine is required only when azureopenai is the ai backend
//if config.Spec.AI.Engine != "" && config.Spec.AI.Backend == v1alpha1.AzureOpenAI {
//	engine := corev1.EnvVar{
//		Name:  "K8SGPT_ENGINE",
//		Value: config.Spec.AI.Engine,
//	}
//	deployment.Spec.Template.Spec.Containers[0].Env = append(
//		deployment.Spec.Template.Spec.Containers[0].Env, engine,
//	)
//} else if config.Spec.AI.Engine != "" && config.Spec.AI.Backend != v1alpha1.AzureOpenAI {
//	return &appsv1.Deployment{}, err.New("Engine is supported only by azureopenai provider.")
//}
//	return &deployment, nil
//}
