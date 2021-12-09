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

package kube

import (
	"os"

	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Clients returns new Kubernetes and Dapr clients
func Clients() (*kubernetes.Clientset, scheme.Interface, error) {
	var config *rest.Config
	var err error
	pathToKubeConfig := os.Getenv("DAPR_DASHBOARD_KUBECONFIG")
	if pathToKubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", pathToKubeConfig)
		if err != nil {
			return nil, nil, err
		}
	}

	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, nil, err
		}
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return kubeClient, nil, err
	}

	daprClient, err := scheme.NewForConfig(config)
	return kubeClient, daprClient, err
}
