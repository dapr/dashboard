/*
Copyright 2021 The Dapr Authors
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

package version

import (
	"errors"
	"github.com/dapr/cli/pkg/kubernetes"
	"github.com/dapr/cli/pkg/standalone"
	"github.com/dapr/dashboard/pkg/kube"
	"strings"
)

// version is the current Dapr dashboard version
var version = "edge"

const operatorName = "dapr-operator"

// GetVersion returns the current dashboard version
func GetVersion() string {
	return version
}

// GetRuntimeVersion returns the current runtime version
func GetRuntimeVersion() (string, error) {
	kubeClient, _, _ := kube.Clients()
	if kubeClient != nil {
		// kubernetes
		sc, err := kubernetes.NewStatusClient()
		if err != nil {
			return "", err
		}

		status, err := sc.Status()
		if err != nil {
			return "", err
		}

		if len(status) == 0 {
			return "", errors.New("Dapr is not installed in your cluster")
		}

		var daprVersion string
		for _, s := range status {
			if s.Name == operatorName {
				daprVersion = s.Version
			}
		}
		return daprVersion, nil
	} else {
		// standalone
		return strings.ReplaceAll(standalone.GetRuntimeVersion(), "\n", ""), nil
	}
}
