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

// ResultSpec defines the desired state of Result

type Event struct {
	Type    string `json:"Type"`
	Reason  string `json:"Reason"`
	Count   int16  `json:"Count"`
	Message string `json:"Message"`
}

type ResultSpec struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Kind      string            `json:"kind"`
	Event     []Event           `json:"Event"`
	Images    []string          `json:"images"`
	Labels    map[string]string `json:"labels"`
}

// ResultStatus defines the observed state of Result
type ResultStatus struct {
	Webhook   string `json:"webhook,omitempty"`
	LifeCycle string `json:"lifecycle,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Result is the Schema for the results API
type Result struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResultSpec   `json:"spec,omitempty"`
	Status ResultStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResultList contains a list of Result
type ResultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Result `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Result{}, &ResultList{})
}
