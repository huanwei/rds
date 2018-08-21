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

package constants

// DefaultMySQLAgentHeathcheckPort is the port on which the mysql-agent's
// healthcheck service runs on.
const DefaultMySQLAgentHeathcheckPort int32 = 10512

// ClusterLabel is applied to all components of a MySQL cluster
const ClusterLabel = "v1alpha1.rds.huanwei.io/cluster"

// MySQLOperatorVersionLabel denotes the version of the MySQLOperator and
// MySQLOperatorAgent running in the cluster.
const MySQLOperatorVersionLabel = "v1alpha1.rds.huanwei.io/version"

// LabelClusterRole specifies the role of a Pod within a Cluster.
const LabelClusterRole = "v1alpha1.rds.huanwei.io/role"

// ClusterRolePrimary denotes a primary InnoDB cluster member.
const ClusterRolePrimary = "primary"

// ClusterRoleSecondary denotes a secondary InnoDB cluster member.
const ClusterRoleSecondary = "secondary"
