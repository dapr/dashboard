package kube

import (
	"os"

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
	daprClient, err := scheme.NewForConfig(config)
	return kubeClient, daprClient, err
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
