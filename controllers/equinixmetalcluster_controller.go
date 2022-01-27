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

package controllers

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-equinixmetal/api/v1beta1"
)

// EquinixMetalClusterReconciler reconciles a EquinixMetalCluster object.
type EquinixMetalClusterReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
}

var ErrNotImplemented = errors.New("not implemented yet")

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EquinixMetalCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *EquinixMetalClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("In Reconcile")

	// TODO(user): your logic here

	return ctrl.Result{}, ErrNotImplemented
}

// SetupWithManager sets up the controller with the Manager.
func (r *EquinixMetalClusterReconciler) SetupWithManager(
	ctx context.Context,
	mgr ctrl.Manager,
	options controller.Options,
) error {
	log := ctrl.LoggerFrom(ctx)

	if r.Client == nil {
		r.Client = mgr.GetClient()
	}

	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("equinixmetalcluster-controller")
	}

	ctrlBuilder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(new(infrav1.EquinixMetalCluster)).
		// Filter out any paused or filtered resources
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		// Filter out any externally managed resources
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		// Watch for changes in CAPI Cluster resources
		Watches(
			&source.Kind{Type: new(clusterv1.Cluster)},
			handler.EnqueueRequestsFromMapFunc(
				util.ClusterToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("EquinixMetalCluster")),
			),
			builder.WithPredicates(predicates.ClusterUnpaused(log)),
		)

	if err := ctrlBuilder.Complete(r); err != nil {
		return fmt.Errorf("failed to create EquinixMetalCluster controller: %w", err)
	}

	return nil
}
