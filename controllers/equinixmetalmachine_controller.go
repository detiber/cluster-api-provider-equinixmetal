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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

// EquinixMetalMachineReconciler reconciles a EquinixMetalMachine object.
type EquinixMetalMachineReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalmachines/finalizers,verbs=update
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=equinixmetalclusters;equinixmetalclusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EquinixMetalMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *EquinixMetalMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("In Reconcile")

	// TODO(user): your logic here

	return ctrl.Result{}, ErrNotImplemented
}

// SetupWithManager sets up the controller with the Manager.
func (r *EquinixMetalMachineReconciler) SetupWithManager(
	ctx context.Context,
	mgr ctrl.Manager,
	options controller.Options,
) error {
	log := ctrl.LoggerFrom(ctx)

	if r.Client == nil {
		r.Client = mgr.GetClient()
	}

	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("equinixmetalmachine-controller")
	}

	equinixMetalMachineMapper, err := util.ClusterToObjectsMapper(
		r.Client,
		new(infrav1.EquinixMetalMachineList),
		mgr.GetScheme(),
	)
	if err != nil {
		return fmt.Errorf("failed to create mapper for Cluster to EquinixMetalMachines: %w", err)
	}

	ctrlBuilder := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(new(infrav1.EquinixMetalMachine)).
		// Filter out any paused or filtered resources
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		// Watch for changes in CAPI Machine resources
		Watches(
			&source.Kind{Type: new(clusterv1.Machine)},
			handler.EnqueueRequestsFromMapFunc(
				util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("EquinixMetalMachine")),
			),
		).
		// Watch for changes in EquinixMetalCluster
		Watches(
			&source.Kind{Type: new(infrav1.EquinixMetalCluster)},
			handler.EnqueueRequestsFromMapFunc(r.equinixMetalClusterToEquinixMetalMachinesMapper(ctx)),
		).
		// Watch for changes in CAPI Cluster resources
		Watches(
			&source.Kind{Type: new(clusterv1.Cluster)},
			handler.EnqueueRequestsFromMapFunc(equinixMetalMachineMapper),
			builder.WithPredicates(predicates.ClusterUnpausedAndInfrastructureReady(log)),
		)

	if err := ctrlBuilder.Complete(r); err != nil {
		return fmt.Errorf("failed to create EquinixMetalMachine controller: %w", err)
	}

	return nil
}

func (r *EquinixMetalMachineReconciler) equinixMetalClusterToEquinixMetalMachinesMapper(
	ctx context.Context,
) handler.MapFunc {
	return func(o client.Object) []ctrl.Request {
		equinixMetalCluster, ok := o.(*infrav1.EquinixMetalCluster)
		if !ok {
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, equinixMetalCluster.ObjectMeta)

		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			return nil
		case err != nil:
			return nil
		}

		labels := map[string]string{clusterv1.ClusterLabelName: cluster.Name}

		machineList := new(clusterv1.MachineList)
		if err := r.Client.List(
			ctx,
			machineList,
			client.InNamespace(equinixMetalCluster.Namespace),
			client.MatchingLabels(labels),
		); err != nil {
			return nil
		}

		result := make([]ctrl.Request, 0, len(machineList.Items))

		for _, m := range machineList.Items {
			if m.Spec.InfrastructureRef.Name == "" {
				continue
			}

			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Name}

			result = append(result, ctrl.Request{NamespacedName: name})
		}

		return result
	}
}
