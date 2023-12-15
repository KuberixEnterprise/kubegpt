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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type AISpec struct {
	// +kubebuilder: default:=true
	Enabled bool `json:"enabled"`
	// +kubebuilder: default:=openai
	Backend string `json:"backend"`
	// +kubebuilder: default:=kr
	Language string `json:"language"`
	// +kubebuilder:default:=gpt-gpt-4-1106-preview
	Model  string     `json:"model"`
	Secret *SecretRef `json:"secret"`
}
type SecretRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type WebhookRef struct {
	// +kubebuilder:validation:Enum=slack
	Type     string `json:"type,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}
type TimerRef struct {
	// +kubebuilder: default:=30
	ErrorInterval int64 `json:"errorInterval"`
	// +kubebuilder: default:=30
	SlackInterval int64 `json:"slackInterval"`
}

type CacheRef struct {
	// +kubebuilder: default:=true
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type KubegptSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Version string  `json:"version"`
	AI *AISpec `json:"ai"`
	//Timer *TimerRef   `json:"timer"`
	Sink  *WebhookRef `json:"sink"`
	Cache *CacheRef   `json:"cache"`
	Timer *TimerRef   `json:"timer"`
}

// KubegptStatus defines the observed state of Kubegpt
type KubegptStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Kubegpt is the Schema for the kubegpts API
type Kubegpt struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubegptSpec   `json:"spec,omitempty"`
	Status KubegptStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KubegptList contains a list of Kubegpt
type KubegptList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kubegpt `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kubegpt{}, &KubegptList{})
}
