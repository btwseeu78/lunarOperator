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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SolarSpec defines the desired state of Solar
type SolarSpec struct {
	MyName string `json:"myName"`
}

// SolarStatus defines the observed state of Solar
type SolarStatus struct {
	MoonName string `json:"moonName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".spec.status.moonName",name="MoonName",type="strin"

// Solar is the Schema for the solars API
type Solar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SolarSpec   `json:"spec,omitempty"`
	Status SolarStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SolarList contains a list of Solar
type SolarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Solar `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Solar{}, &SolarList{})
}
