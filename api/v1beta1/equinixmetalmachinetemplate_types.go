/*
Copyright 2022.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EquinixMetalMachineTemplateResource describes the data needed to create am EquinixMetalMachine from a template.
type EquinixMetalMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec EquinixMetalMachineSpec `json:"spec"`
}

// EquinixMetalMachineTemplateSpec defines the desired state of EquinixMetalMachineTemplate.
type EquinixMetalMachineTemplateSpec struct {
	Template EquinixMetalMachineTemplateResource `json:"template"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=equinixmetalmachinetemplates,scope=Namespaced,categories=cluster-api
//+kubebuilder:storageversion

// EquinixMetalMachineTemplate is the Schema for the equinixmetalmachinetemplates API.
type EquinixMetalMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"` //nolint:tagliatelle
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EquinixMetalMachineTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// EquinixMetalMachineTemplateList contains a list of EquinixMetalMachineTemplate.
type EquinixMetalMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"` //nolint:tagliatelle
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EquinixMetalMachineTemplate `json:"items"`
}

func init() { //nolint:gochecknoinits
	SchemeBuilder.Register(new(EquinixMetalMachineTemplate), new(EquinixMetalMachineTemplateList))
}
