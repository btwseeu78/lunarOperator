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

// MoonSpec defines the desired state of Moon
type MoonSpec struct {
	//+kubebuilder:default=1
	//+kubebuilder:validation:Minimum=0
	//+optional
	Replicas *int32 `json:"replicas"`

	//+optional
	SunName string `json:"sunName,omitempty"`
}

// MoonStatus defines the observed state of Moon
type MoonStatus struct {
	FoundSun   string       `json:"foundSun"`
	LastSynced *metav1.Time `json:"lastSynced"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.foundSun",name="Sun",type="string"
//+kubebuilder:printcolumn:JSONPath=".spec.replicas",name="Desired",type="integer"
//+kubebuilder:printcolumn:JSONPath=".status.lastSynced",name="LastSynced",type="string"

// Moon is the Schema for the moons API
type Moon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MoonSpec   `json:"spec,omitempty"`
	Status MoonStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MoonList contains a list of Moon
type MoonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Moon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Moon{}, &MoonList{})
}
