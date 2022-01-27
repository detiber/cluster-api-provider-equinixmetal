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

package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-equinixmetal/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-equinixmetal/controllers"
)

const (
	defaultEqunixMetalClusterConcurrency  = 5
	defaultEquinixMetalMachineConcurrency = 10
	defaultWebhookPort                    = 9443
	defaultSyncPeriod                     = 10 * time.Minute
)

type config struct {
	metricsBindAddr                string
	enableLeaderElection           bool
	leaderElectionNamespace        string
	watchNamespace                 string
	watchFilterValue               string
	profilerAddress                string
	equinixMetalClusterConcurrency int
	equinixMetalMachineConcurrency int
	syncPeriod                     time.Duration
	webhookPort                    int
	webhookCertDir                 string
	healthAddr                     string
}

func main() { //nolint:funlen
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())

	config := new(config)

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme))
	utilruntime.Must(infrav1beta1.AddToScheme(scheme))

	configureFlags(pflag.CommandLine, config)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klogr.New())
	setupLog := ctrl.Log.WithName("setup")

	if config.watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", config.watchNamespace)
	}

	if config.profilerAddress != "" {
		setupLog.Info("Profiler listening for requests", "profiler-address", config.profilerAddress)

		go func() {
			setupLog.Error(http.ListenAndServe(config.profilerAddress, nil), "listen and serve error")
		}()
	}

	ctx := ctrl.SetupSignalHandler()

	restConfig := ctrl.GetConfigOrDie()
	restConfig.UserAgent = "cluster-api-provider-equinixmetal-controller"

	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{ //nolint:exhaustivestruct
		Scheme:                     scheme,
		MetricsBindAddress:         config.metricsBindAddr,
		LeaderElection:             config.enableLeaderElection,
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaderElectionID:           "controller-leader-elect-capem",
		LeaderElectionNamespace:    config.leaderElectionNamespace,
		SyncPeriod:                 &config.syncPeriod,
		Namespace:                  config.watchNamespace,
		Port:                       config.webhookPort,
		CertDir:                    config.webhookCertDir,
		HealthProbeBindAddress:     config.healthAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("equinixmetal-controller"))

	setupLog.V(1).Info(fmt.Sprintf("feature gates: %+v\n", feature.Gates))

	if err := setupControllers(ctx, mgr, config); err != nil {
		setupLog.Error(err, "failed to configure controllers")
		os.Exit(1)
	}

	if err := setupWebhooks(mgr); err != nil {
		setupLog.Error(err, "failed to configure controllers")
		os.Exit(1)
	}

	if err := setupChecks(mgr); err != nil {
		setupLog.Error(err, "failed to configure health and readiness checks")
		os.Exit(1)
	}

	setupLog.Info("starting manager", "version", version.Get().String())

	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func configureFlags(flagset *pflag.FlagSet, config *config) { //nolint:funlen
	flagset.StringVar(
		&config.metricsBindAddr,
		"metrics-bind-addr",
		"localhost:8080",
		"The address the metric endpoint binds to.",
	)

	flagset.BoolVar(
		&config.enableLeaderElection,
		"leader-elect",
		false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.",
	)

	flagset.StringVar(
		&config.watchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. "+
			"If unspecified, the controller watches for cluster-api objects across all namespaces.",
	)

	flagset.StringVar(
		&config.leaderElectionNamespace,
		"leader-elect-namespace",
		"",
		"Namespace that the controller performs leader election in. "+
			"If unspecified, the controller will discover which namespace it is running in.",
	)

	flagset.StringVar(
		&config.profilerAddress,
		"profiler-address",
		"",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)",
	)

	flagset.IntVar(&config.equinixMetalClusterConcurrency,
		"equinixmetalcluster-concurrency",
		defaultEqunixMetalClusterConcurrency,
		"Number of EquinixMetalClusters to process simultaneously",
	)

	flagset.IntVar(&config.equinixMetalMachineConcurrency,
		"equinixmetalmachine-concurrency",
		defaultEquinixMetalMachineConcurrency,
		"Number of EquinixMetalMachines to process simultaneously",
	)

	flagset.DurationVar(&config.syncPeriod,
		"sync-period",
		defaultSyncPeriod,
		"The minimum interval at which watched resources are reconciled",
	)

	flagset.IntVar(&config.webhookPort,
		"webhook-port",
		defaultWebhookPort,
		"Webhook Server port.",
	)

	flagset.StringVar(&config.webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")

	flagset.StringVar(&config.healthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)

	flagset.StringVar(
		&config.watchFilterValue,
		"watch-filter",
		"",
		fmt.Sprintf(
			"Label value that the controller watches to reconcile cluster-api objects. "+
				"Label key is always %s. If unspecified, the controller watches for all cluster-api objects.",
			clusterv1.WatchLabel,
		),
	)

	feature.MutableGates.AddFlag(flagset)
}

func setupControllers(ctx context.Context, mgr ctrl.Manager, config *config) error {
	if err := (&controllers.EquinixMetalClusterReconciler{ //nolint:exhaustivestruct
		WatchFilterValue: config.watchFilterValue,
	}).SetupWithManager(
		ctx,
		mgr,
		controller.Options{ //nolint:exhaustivestruct
			MaxConcurrentReconciles: config.equinixMetalClusterConcurrency,
			RecoverPanic:            true,
		},
	); err != nil {
		return fmt.Errorf("unable to create EquinixMetalCluster controller: %w", err)
	}

	if err := (&controllers.EquinixMetalMachineReconciler{ //nolint:exhaustivestruct
		WatchFilterValue: config.watchFilterValue,
	}).SetupWithManager(
		ctx,
		mgr,
		controller.Options{ //nolint:exhaustivestruct
			MaxConcurrentReconciles: config.equinixMetalMachineConcurrency,
			RecoverPanic:            true,
		},
	); err != nil {
		return fmt.Errorf("unable to create EquinixMetalMachine controller: %w", err)
	}

	return nil
}

func setupWebhooks(mgr ctrl.Manager) error {
	if err := new(infrav1beta1.EquinixMetalCluster).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create EquinixMetalCluster webhook: %w", err)
	}

	if err := new(infrav1beta1.EquinixMetalMachine).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create EquinixMetalMachine webhook: %w", err)
	}

	if err := new(infrav1beta1.EquinixMetalMachineTemplate).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create EquinixMetalMachineTemplate webhook: %w", err)
	}

	return nil
}

func setupChecks(mgr ctrl.Manager) error {
	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		return fmt.Errorf("unable to create readiness check: %w", err)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		return fmt.Errorf("unable to create health check: %w", err)
	}

	return nil
}
