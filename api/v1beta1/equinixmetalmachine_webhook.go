/*
Copyright 2021 The Kubernetes Authors.

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
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func (m *EquinixMetalMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	if err := ctrl.NewWebhookManagedBy(mgr).For(m).Complete(); err != nil {
		return fmt.Errorf("failed to create EquinixMetalMachine webhook: %w", err)
	}

	return nil
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-equinixmetalmachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachines,versions=v1beta1,name=validation.equinixmetalmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-equinixmetalmachine,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachines,versions=v1beta1,name=default.equinixmetalmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (m *EquinixMetalMachine) ValidateCreate() error {
	machineLog := logf.Log.WithName("equinixmetalmachine-resource")
	machineLog.Info("validate create", "name", m.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (m *EquinixMetalMachine) ValidateUpdate(old runtime.Object) error {
	machineLog := logf.Log.WithName("equinixmetalmachine-resource")
	machineLog.Info("validate update", "name", m.Name)

	newEquinixMetalMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(m)
	if err != nil {
		return apierrors.NewInvalid(GroupVersion.WithKind("EquinixMetalMachine").GroupKind(), m.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert new EquinixMetalMachine to unstructured object")),
		})
	}

	oldEquinixMetalMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(old)
	if err != nil {
		return apierrors.NewInvalid(GroupVersion.WithKind("EquinixMetalMachine").GroupKind(), m.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert old EquinixMetalMachine to unstructured object")),
		})
	}

	newEquinixMetalMachineSpec, _ := newEquinixMetalMachine["spec"].(map[string]interface{})
	oldEquinixMetalMachineSpec, _ := oldEquinixMetalMachine["spec"].(map[string]interface{})

	// allow changes to providerID
	delete(oldEquinixMetalMachineSpec, "providerID")
	delete(newEquinixMetalMachineSpec, "providerID")

	// allow changes to tags
	delete(oldEquinixMetalMachineSpec, "tags")
	delete(newEquinixMetalMachineSpec, "tags")

	if !reflect.DeepEqual(oldEquinixMetalMachineSpec, newEquinixMetalMachineSpec) {
		return apierrors.NewInvalid(GroupVersion.WithKind("EquinixMetalMachine").GroupKind(), m.Name, field.ErrorList{
			field.Forbidden(field.NewPath("spec"), "cannot be modified"),
		})
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (m *EquinixMetalMachine) ValidateDelete() error {
	machineLog := logf.Log.WithName("equinixmetalmachine-resource")
	machineLog.Info("validate delete", "name", m.Name)

	return nil
}

// Default implements webhookutil.defaulter so a webhook will be registered for the type.
func (m *EquinixMetalMachine) Default() {
	machineLog := logf.Log.WithName("equinixmetalmachine-resource")
	machineLog.Info("default", "name", m.Name)
}
