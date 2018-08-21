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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MinimumMySQLVersion is the minimum version of MySQL server supported by the
// MySQL Operator.
const MinimumMySQLVersion = "8.0.11"

// ClusterSpec defines the attributes a user can specify when creating a cluster
type ClusterSpec struct {
	// Version defines the MySQL Docker image version.
	Version string `json:"version"`
	// Members defines the number of MySQL instances in a cluster
	Members int32 `json:"members,omitempty"`
	// BaseServerID defines the base number used to create unique server_id
	// for MySQL instances in the cluster. Valid range 1 to 4294967286.
	// If omitted in the manifest file (or set to 0) defaultBaseServerID
	// value will be used.
	BaseServerID uint32 `json:"baseServerId,omitempty"`
	// MultiMaster defines the mode of the MySQL cluster. If set to true,
	// all instances will be R/W. If false (the default), only a single instance
	// will be R/W and the rest will be R/O.
	MultiMaster bool `json:"multiMaster,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, affinity will define the pod's scheduling constraints
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// VolumeClaimTemplate allows a user to specify how volumes inside a MySQL cluster
	// +optional
	VolumeClaimTemplate *corev1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
	// BackupVolumeClaimTemplate allows a user to specify a volume to temporarily store the
	// data for a backup prior to it being shipped to object storage.
	// +optional
	BackupVolumeClaimTemplate *corev1.PersistentVolumeClaim `json:"backupVolumeClaimTemplate,omitempty"`
	// If defined, we use this secret for configuring the MYSQL_ROOT_PASSWORD
	// If it is not set we generate a secret dynamically
	// +optional
	RootPasswordSecret *corev1.LocalObjectReference `json:"rootPasswordSecret,omitempty"`
	// Config allows a user to specify a custom configuration file for MySQL.
	// +optional
	Config *corev1.LocalObjectReference `json:"config,omitempty"`
	// SSLSecret allows a user to specify custom CA certificate, server certificate
	// and server key for group replication SSL.
	// +optional
	SSLSecret *corev1.LocalObjectReference `json:"sslSecret,omitempty"`
}

// ClusterConditionType represents a valid condition of a Cluster.
type ClusterConditionType string

const (
	// ClusterReady means the Cluster is able to service requests.
	ClusterReady ClusterConditionType = "Ready"
)

// ClusterCondition describes the observed state of a Cluster at a certain point.
type ClusterCondition struct {
	Type   ClusterConditionType
	Status corev1.ConditionStatus
	// +optional
	LastTransitionTime metav1.Time
	// +optional
	Reason string
	// +optional
	Message string
}

// ClusterStatus defines the current status of a MySQL cluster
// propagating useful information back to the cluster admin
type ClusterStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// +optional
	Conditions []ClusterCondition
}

// +genclient
// +genclient:noStatus
// +resourceName=mysqlclusters
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster represents a cluster spec and associated metadata
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ClusterSpec   `json:"spec"`
	Status ClusterStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList is a placeholder type for a list of MySQL clusters
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Cluster `json:"items"`
}
