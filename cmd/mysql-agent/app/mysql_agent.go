/*
Copyright The Kubernetes Authors.

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

package app

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"

	kubeinformers "k8s.io/client-go/informers"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"

	cluster "github.com/huanwei/rds/pkg/mysql-agent/cluster"
	clustermgr "github.com/huanwei/rds/pkg/mysql-agent/controllers/cluster/manager"
	agentopts "github.com/huanwei/rds/pkg/mysql-agent/options"
	clientset "github.com/huanwei/rds/pkg/mysql-operator/generated/clientset/versioned"
	informers "github.com/huanwei/rds/pkg/mysql-operator/generated/informers/externalversions"
	signals "github.com/huanwei/rds/pkg/util/signals"
)

// resyncPeriod computes the time interval a shared informer waits before
// resyncing with the api server.
func resyncPeriod(opts *agentopts.MySQLAgentOpts) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(opts.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run runs the MySQL backup controller. It should never exit.
func Run(opts *agentopts.MySQLAgentOpts) error {
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Set up signals so we handle the first shutdown signal gracefully.
	signals.SetupSignalHandler(cancelFunc)

	// Set up healthchecks (liveness and readiness).
	checkInCluster, err := cluster.NewHealthCheck()
	if err != nil {
		glog.Fatal(err)
	}
	health := healthcheck.NewHandler()
	health.AddReadinessCheck("node-in-cluster", checkInCluster)
	go func() {
		glog.Fatal(http.ListenAndServe(
			net.JoinHostPort(opts.Address, strconv.Itoa(int(opts.HealthcheckPort))),
			health,
		))
	}()

	kubeclient := kubernetes.NewForConfigOrDie(kubeconfig)
	mysqlopClient := clientset.NewForConfigOrDie(kubeconfig)

	kubeInformerFactory := kubeinformers.NewFilteredSharedInformerFactory(kubeclient, resyncPeriod(opts)(), opts.Namespace, nil)
	sharedInformerFactory := informers.NewFilteredSharedInformerFactory(mysqlopClient, 0, opts.Namespace, nil)

	var wg sync.WaitGroup

	manager, err := clustermgr.NewLocalClusterManger(kubeclient, kubeInformerFactory)
	if err != nil {
		return errors.Wrap(err, "failed to create new local MySQL InnoDB cluster manager")
	}

	// Block until local instance successfully initialised.
	for !manager.Sync(ctx) {
		time.Sleep(10 * time.Second)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		manager.Run(ctx)
	}()

	// Shared informers have to be started after ALL controllers.
	go sharedInformerFactory.Start(ctx.Done())
	go kubeInformerFactory.Start(ctx.Done())

	<-ctx.Done()

	glog.Info("Waiting for all controllers to shut down gracefully")
	wg.Wait()

	return nil
}
