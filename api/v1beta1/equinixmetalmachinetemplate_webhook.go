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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func (m *EquinixMetalMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	if err := ctrl.NewWebhookManagedBy(mgr).For(m).Complete(); err != nil {
		return fmt.Errorf("failed to create EquinixMetalMachineTemplate webhook: %w", err)
	}

	return nil
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-equinixmetalmachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachinetemplates,versions=v1beta1,name=validation.equinixmetalmachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-equinixmetalmachinetemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachinetemplates,versions=v1beta1,name=default.equinixmetalmachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (m *EquinixMetalMachineTemplate) ValidateCreate() error {
	machineTemplateLog := logf.Log.WithName("equinixmetalmachinetemplate-resource")
	machineTemplateLog.Info("validate create", "name", m.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (m *EquinixMetalMachineTemplate) ValidateUpdate(old runtime.Object) error {
	machineTemplateLog := logf.Log.WithName("equinixmetalmachinetemplate-resource")
	machineTemplateLog.Info("validate update", "name", m.Name)

	oldEquinixMetalMachineTemplate, _ := old.(*EquinixMetalMachineTemplate)

	if !reflect.DeepEqual(m.Spec, oldEquinixMetalMachineTemplate.Spec) {
		return apierrors.NewBadRequest("EquinixMetalMachineTemplate.Spec is immutable")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (m *EquinixMetalMachineTemplate) ValidateDelete() error {
	machineTemplateLog := logf.Log.WithName("equinixmetalmachinetemplate-resource")
	machineTemplateLog.Info("validate delete", "name", m.Name)

	return nil
}

// Default implements webhookutil.defaulter so a webhook will be registered for the type.
func (m *EquinixMetalMachineTemplate) Default() {
	machineTemplateLog := logf.Log.WithName("equinixmetalmachinetemplate-resource")
	machineTemplateLog.Info("default", "name", m.Name)
}
