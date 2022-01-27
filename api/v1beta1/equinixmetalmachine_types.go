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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows ReconcileEquinixMetalMachine to clean up EquinixMetal resources before
	// removing it from the apiserver.
	MachineFinalizer = "equinixmetalmachine.infrastructure.cluster.x-k8s.io"
)

const (
	// DeviceReadyCondition reports on current status of the Equinix Metal device.
	// Ready indicates the instance is in a Running state.
	DeviceReadyCondition clusterv1.ConditionType = "InstanceReady"

	// InstanceNotFoundReason used when the instance couldn't be retrieved.
	InstanceNotFoundReason = "InstanceNotFound"
	// InstanceTerminatedReason instance is in a terminated state.
	InstanceTerminatedReason = "InstanceTerminated"
	// InstanceStoppedReason instance is in a stopped state.
	InstanceStoppedReason = "InstanceStopped"
	// InstanceNotReadyReason used when the instance is in a pending state.
	InstanceNotReadyReason = "InstanceNotReady"
	// InstanceProvisionStartedReason set when the provisioning of an instance started.
	InstanceProvisionStartedReason = "InstanceProvisionStarted"
	// InstanceProvisionFailedReason used for failures during instance provisioning.
	InstanceProvisionFailedReason = "InstanceProvisionFailed"
	// WaitingForClusterInfrastructureReason used when machine is waiting for cluster infrastructure to be ready
	// before proceeding.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"
	// WaitingForBootstrapDataReason used when machine is waiting for bootstrap data to be ready before proceeding.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"
)

// EquinixMetalResourceStatus describes the status of a EquinixMetal resource.
type EquinixMetalResourceStatus string

const (
	// EquinixMetalResourceStatusNew represents a EquinixMetal resource requested.
	// The EquinixMetal infrastucture uses a queue to avoid any abuse. So a resource
	// does not get created straigh away but it can wait for a bit in a queue.
	EquinixMetalResourceStatusNew = EquinixMetalResourceStatus("new")
	// EquinixMetalResourceStatusQueued represents a device waiting for his turn to be provisioned.
	// Time in queue depends on how many creation requests you already issued, or
	// from how many resources waiting to be deleted we have for you.
	EquinixMetalResourceStatusQueued = EquinixMetalResourceStatus("queued")
	// EquinixMetalResourceStatusProvisioning represents a resource that got dequeued
	// and it is actively processed by a worker.
	EquinixMetalResourceStatusProvisioning = EquinixMetalResourceStatus("provisioning")
	// EquinixMetalResourceStatusRunning represents a EquinixMetal resource already provisioned and in a active state.
	EquinixMetalResourceStatusRunning = EquinixMetalResourceStatus("active")
	// EquinixMetalResourceStatusErrored represents a EquinixMetal resource in a errored state.
	EquinixMetalResourceStatusErrored = EquinixMetalResourceStatus("errored")
	// EquinixMetalResourceStatusOff represents a EquinixMetal resource in off state.
	EquinixMetalResourceStatusOff = EquinixMetalResourceStatus("off")
)

// EquinixMetalMachineSpec defines the desired state of EquinixMetalMachine.
type EquinixMetalMachineSpec struct {
	OS           string   `json:"os"`
	BillingCycle string   `json:"billingCycle"`
	MachineType  string   `json:"machineType"`
	SSHKeys      []string `json:"sshKeys,omitempty"`

	// Metro represents the EquinixMetal metro for this machine.
	// Override from the EquinixMetalCluster spec.
	// +optional
	Metro string `json:"metro,omitempty"`

	// Facility represents the EquinixMetal facility for this machine.
	// Override from the EquinixMetalCluster spec.
	// +optional
	Facility string `json:"facility,omitempty"`

	// IPXEUrl can be used to set the pxe boot url when using custom OSes with this provider.
	// Note that OS should also be set to "custom_ipxe" if using this value.
	// +optional
	IPXEUrl string `json:"ipxeURL,omitempty"`

	// HardwareReservationID is the unique device hardware reservation ID, a comma separated list of
	// hardware reservation IDs, or `next-available` to automatically let the EquinixMetal api determine one.
	// +optional
	HardwareReservationID string `json:"hardwareReservationID,omitempty"`

	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// Tags is an optional set of tags to add to EquinixMetal resources managed by the EquinixMetal provider.
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// EquinixMetalMachineStatus defines the observed state of EquinixMetalMachine.
type EquinixMetalMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the EquinixMetal device associated addresses.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// InstanceStatus is the status of the EquinixMetal device instance for this machine.
	// +optional
	InstanceStatus *EquinixMetalResourceStatus `json:"instanceStatus,omitempty"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *capierrors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// Conditions defines current service state of the EquinixMetalMachine.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=equinixmetalmachines,scope=Namespaced,categories=cluster-api
//+kubebuilder:storageversion
//+kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this EquinixMetalMachine belongs"
//+kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="EquinixMetal instance state"
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
//+kubebuilder:printcolumn:name="InstanceID",type="string",JSONPath=".spec.providerID",description="EquinixMetal instance ID"
//+kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this EquinixMetalMachine"

// EquinixMetalMachine is the Schema for the equinixmetalmachines API.
type EquinixMetalMachine struct {
	metav1.TypeMeta   `json:",inline"` //nolint:tagliatelle
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EquinixMetalMachineSpec   `json:"spec,omitempty"`
	Status EquinixMetalMachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EquinixMetalMachineList contains a list of EquinixMetalMachine.
type EquinixMetalMachineList struct {
	metav1.TypeMeta `json:",inline"` //nolint:tagliatelle
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EquinixMetalMachine `json:"items"`
}

// GetConditions returns the list of conditions for an EquinixMetalMachine API object.
func (m *EquinixMetalMachine) GetConditions() clusterv1.Conditions {
	return m.Status.Conditions
}

// SetConditions will set the given conditions on an EquinixMetalMachine object.
func (m *EquinixMetalMachine) SetConditions(conditions clusterv1.Conditions) {
	m.Status.Conditions = conditions
}

func init() { //nolint:gochecknoinits
	SchemeBuilder.Register(new(EquinixMetalMachine), new(EquinixMetalMachineList))
}
