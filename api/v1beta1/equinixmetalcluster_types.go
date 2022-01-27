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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// NetworkInfrastructureReadyCondition reports of current status of cluster infrastructure.
	NetworkInfrastructureReadyCondition clusterv1.ConditionType = "NetworkInfrastructureReady"
)

// EquinixMetalClusterSpec defines the desired state of EquinixMetalCluster.
type EquinixMetalClusterSpec struct {
	// ProjectID represents the Equinix Metal Project where this cluster will be placed into.
	ProjectID string `json:"projectID"`

	// Metro represents the Equinix Metal metro for this cluster.
	Metro string `json:"metro,omitempty"`

	// Facility represents the Equinix Metal facility for this cluster.
	Facility string `json:"facility,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// EquinixMetalClusterStatus defines the observed state of EquinixMetalCluster.
type EquinixMetalClusterStatus struct {
	// Ready denotes that the cluster (infrastructure) is ready.
	// +optional
	Ready bool `json:"ready"`

	// FailureDomains specifies the list of unique failure domains for the location/region of the cluster.
	// A FailureDomain maps to a Facility with an EquinixMetal Metro. An
	// This list will be used by Cluster API to try and spread the machines across the failure domains.
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// Conditions defines current service state of the EquinixMetalCluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=equinixmetalclusters,scope=Namespaced,categories=cluster-api
//+kubebuilder:storageversion
//+kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this EquinixMetalCluster belongs"
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="EquinixMetalCluster ready status"

// EquinixMetalCluster is the Schema for the equinixmetalclusters API.
type EquinixMetalCluster struct {
	metav1.TypeMeta   `json:",inline"` //nolint:tagliatelle
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EquinixMetalClusterSpec   `json:"spec,omitempty"`
	Status EquinixMetalClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EquinixMetalClusterList contains a list of EquinixMetalCluster.
type EquinixMetalClusterList struct {
	metav1.TypeMeta `json:",inline"` //nolint:tagliatelle
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EquinixMetalCluster `json:"items"`
}

// GetConditions returns the list of conditions for an EquinixMetalCluster API object.
func (c *EquinixMetalCluster) GetConditions() clusterv1.Conditions {
	return c.Status.Conditions
}

// SetConditions will set the given conditions on an EquinixMetalCluster object.
func (c *EquinixMetalCluster) SetConditions(conditions clusterv1.Conditions) {
	c.Status.Conditions = conditions
}

func init() { //nolint:gochecknoinits
	SchemeBuilder.Register(new(EquinixMetalCluster), new(EquinixMetalClusterList))
}
