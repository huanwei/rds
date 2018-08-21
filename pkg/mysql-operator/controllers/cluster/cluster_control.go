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

package cluster

import (
	"fmt"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/util/retry"

	v1alpha1 "github.com/huanwei/rds/pkg/mysql-operator/apis/rds/v1alpha1"
	clientset "github.com/huanwei/rds/pkg/mysql-operator/generated/clientset/versioned"
	listersv1alpha1 "github.com/huanwei/rds/pkg/mysql-operator/generated/listers/rds/v1alpha1"
)

type clusterUpdaterInterface interface {
	UpdateClusterStatus(cluster *v1alpha1.Cluster, status *v1alpha1.ClusterStatus) error
	UpdateClusterLabels(cluster *v1alpha1.Cluster, lbls labels.Set) error
}

type clusterUpdater struct {
	client clientset.Interface
	lister listersv1alpha1.ClusterLister
}

func newClusterUpdater(client clientset.Interface, lister listersv1alpha1.ClusterLister) clusterUpdaterInterface {
	return &clusterUpdater{client: client, lister: lister}
}

func (cu *clusterUpdater) UpdateClusterStatus(cluster *v1alpha1.Cluster, status *v1alpha1.ClusterStatus) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		cluster.Status = *status
		_, updateErr := cu.client.RDSV1alpha1().Clusters(cluster.Namespace).Update(cluster)
		if updateErr == nil {
			return nil
		}

		updated, err := cu.lister.Clusters(cluster.Namespace).Get(cluster.Name)
		if err != nil {
			glog.Errorf("Error getting updated Cluster %s/%s: %v", cluster.Namespace, cluster.Name, err)
			return err
		}

		// Copy the Cluster so we don't mutate the cache.
		cluster = updated.DeepCopy()
		return updateErr
	})
}

func (cu *clusterUpdater) UpdateClusterLabels(cluster *v1alpha1.Cluster, lbls labels.Set) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		cluster.Labels = labels.Merge(labels.Set(cluster.Labels), lbls)
		_, updateErr := cu.client.RDSV1alpha1().Clusters(cluster.Namespace).Update(cluster)
		if updateErr == nil {
			return nil
		}

		key := fmt.Sprintf("%s/%s", cluster.GetNamespace(), cluster.GetName())
		glog.V(4).Infof("Conflict updating Cluster labels. Getting updated Cluster %s from cache...", key)

		updated, err := cu.lister.Clusters(cluster.GetNamespace()).Get(cluster.GetName())
		if err != nil {
			glog.Errorf("Error getting updated Cluster %s: %v", key, err)
			return err
		}

		// Copy the Cluster so we don't mutate the cache.
		cluster = updated.DeepCopy()
		return updateErr
	})
}
