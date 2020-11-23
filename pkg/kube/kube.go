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
