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
	"context"
	"github.com/dapr/cli/pkg/standalone"
	"github.com/dapr/dashboard/pkg/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		ctx := context.Background()
		// kubernetes
		options := metav1.ListOptions{}
		deployments, err := kubeClient.AppsV1().Deployments(v1.NamespaceAll).List(ctx, options)
		if err != nil {
			return "", err
		}
		for _, deployment := range deployments.Items {
			if deployment.Name == operatorName {
				image := deployment.Spec.Template.Spec.Containers[0].Image
				daprVersion := image[strings.IndexAny(image, ":")+1:]
				return daprVersion, nil
			}
		}
	} else {
		// standalone
		return strings.ReplaceAll(standalone.GetRuntimeVersion(), "\n", ""), nil
	}
	return "", nil
}
