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
	"fmt"
	"strconv"

	"github.com/coreos/go-semver/semver"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateCluster(c *Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateClusterMetadata(c.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateClusterSpec(c.Spec, field.NewPath("spec"))...)
	allErrs = append(allErrs, validateClusterStatus(c.Status, field.NewPath("status"))...)
	return allErrs
}

func validateClusterMetadata(m metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateName(m.Name, fldPath.Child("name"))...)

	return allErrs
}

func validateName(name string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(name) > ClusterNameMaxLen {
		msg := fmt.Sprintf("longer than maximum supported length %d (see: https://bugs.mysql.com/bug.php?id=90601)", ClusterNameMaxLen)
		allErrs = append(allErrs, field.Invalid(fldPath, name, msg))
	}

	return allErrs
}

func validateClusterSpec(s ClusterSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateVersion(s.Version, fldPath.Child("version"))...)
	allErrs = append(allErrs, validateMembers(s.Members, fldPath.Child("members"))...)
	allErrs = append(allErrs, validateBaseServerID(s.BaseServerID, fldPath.Child("baseServerId"))...)

	return allErrs
}

func validateClusterStatus(s ClusterStatus, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}

func validateVersion(version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	min, err := semver.NewVersion(MinimumMySQLVersion)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, fmt.Errorf("unable to parse minimum MySQL version: %v", err)))
	}

	given, err := semver.NewVersion(version)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, version, fmt.Sprintf("unable to parse MySQL version: %v", err)))
	}

	if len(allErrs) == 0 {
		if given.Compare(*min) == -1 {
			allErrs = append(allErrs, field.Invalid(fldPath, version, fmt.Sprintf("minimum supported MySQL version is %s", MinimumMySQLVersion)))
		}
	}

	return allErrs
}

func validateBaseServerID(baseServerID uint32, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if baseServerID <= maxBaseServerID {
		return allErrs
	}
	return append(allErrs, field.Invalid(fldPath, strconv.FormatUint(uint64(baseServerID), 10), "invalid baseServerId specified"))
}

func validateMembers(members int32, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if members < 1 || members > MaxInnoDBClusterMembers {
		allErrs = append(allErrs, field.Invalid(fldPath, members, "InnoDB clustering supports between 1-9 members"))
	}
	return allErrs
}
