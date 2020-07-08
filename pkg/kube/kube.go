package kube

import (
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Clients returns new Kubernetes and Dapr clients
func Clients() (*kubernetes.Clientset, scheme.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return kubeClient, nil, err
	}

	daprClient, err := scheme.NewForConfig(config)
	return kubeClient, daprClient, err
}
