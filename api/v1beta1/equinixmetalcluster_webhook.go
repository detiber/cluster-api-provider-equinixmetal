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
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (c *EquinixMetalCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	if err := ctrl.NewWebhookManagedBy(mgr).For(c).Complete(); err != nil {
		return fmt.Errorf("failed to create EquinixMetalCluster webhook: %w", err)
	}

	return nil
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-equinixmetalcluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalclusters,versions=v1beta1,name=validation.equinixmetalcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-equinixmetalcluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalclusters,versions=v1beta1,name=default.equinixmetalcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (c *EquinixMetalCluster) Default() {
	clusterlog := logf.Log.WithName("equinixmetalcluster-resource")
	clusterlog.Info("default", "name", c.Name)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (c *EquinixMetalCluster) ValidateCreate() error {
	clusterlog := logf.Log.WithName("equinixmetalcluster-resource")
	clusterlog.Info("validate create", "name", c.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (c *EquinixMetalCluster) ValidateUpdate(oldRaw runtime.Object) error {
	clusterlog := logf.Log.WithName("equinixmetalcluster-resource")
	clusterlog.Info("validate update", "name", c.Name)

	var allErrs field.ErrorList

	old, _ := oldRaw.(*EquinixMetalCluster)

	if !reflect.DeepEqual(c.Spec.ProjectID, old.Spec.ProjectID) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "projectID"),
				c.Spec.ProjectID, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(c.Spec.Facility, old.Spec.Facility) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "Facility"),
				c.Spec.Facility, "field is immutable"),
		)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(GroupVersion.WithKind("EquinixMetalCluster").GroupKind(), c.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (c *EquinixMetalCluster) ValidateDelete() error {
	clusterlog := logf.Log.WithName("equinixmetalcluster-resource")
	clusterlog.Info("validate delete", "name", c.Name)

	return nil
}
