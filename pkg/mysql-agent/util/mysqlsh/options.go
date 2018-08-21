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

package mysqlsh

import (
	"fmt"
	"strings"
)

// Options holds the options passed to individual mysqlsh commands.
type Options map[string]string

// String encodes options as a Python dictionary string.
func (opts Options) String() string {
	vals := []string{}
	for k, v := range opts {
		vals = append(vals, fmt.Sprintf("'%s': %s", k, quoted(v)))
	}
	return fmt.Sprintf("{%s}", strings.Join(vals, ", "))
}

// quoted handles quoting string options vs. not quoting boolean options.
func quoted(s string) string {
	switch strings.ToLower(s) {
	case "true":
		return "True"
	case "false":
		return "False"
	default:
		return fmt.Sprintf("'%s'", s)
	}
}