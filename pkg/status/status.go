package kubernetes

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dapr/dashboard/pkg/age"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubernetes "k8s.io/client-go/kubernetes"
)

// Status is an interface to interact with Dapr control plane information
type Status interface {
	Supported() bool
	Get() []StatusOutput
}

type status struct {
	kubeClient *kubernetes.Clientset
}

var (
	controlPlaneLabels = [...]string{"dapr-operator", "dapr-sentry", "dapr-placement", "dapr-sidecar-injector"}
)

// StatusOutput represents the status of a named Dapr resource.
type StatusOutput struct {
	Name      string `csv:"NAME"`
	Namespace string `csv:"NAMESPACE"`
	Healthy   string `csv:"HEALTHY"`
	Status    string `csv:"STATUS"`
	Version   string `csv:"VERSION"`
	Age       string `csv:"AGE"`
	Created   string `csv:"CREATED"`
}

// NewStatus returns a new Status instance
func NewStatus(kubeClient *kubernetes.Clientset) Status {
	return &status{
		kubeClient: kubeClient,
	}
}

// Supported checks whether or not the kubernetes client is available to access the contol plane statuses
func (s *status) Supported() bool {
	return s.kubeClient != nil
}

// Get lists the status of each of the Dapr control plane services
func (s *status) Get() []StatusOutput {
	var wg sync.WaitGroup
	wg.Add(len(controlPlaneLabels))

	m := sync.Mutex{}
	statuses := []StatusOutput{}

	for _, lbl := range controlPlaneLabels {
		go func(label string) {
			options := v1.ListOptions{}
			labelSelector := map[string]string{
				"app": label,
			}
			options.LabelSelector = labels.FormatLabels(labelSelector)

			p, err := s.kubeClient.CoreV1().Pods(v1.NamespaceAll).List(options)
			if err == nil && len(p.Items) == 1 {
				pod := p.Items[0]

				image := pod.Spec.Containers[0].Image
				namespace := pod.GetNamespace()
				age := age.GetAge(pod.CreationTimestamp.Time)
				created := pod.CreationTimestamp.Format("2006-01-02 15:04.05")
				version := image[strings.IndexAny(image, ":")+1:]
				status := ""

				if pod.Status.ContainerStatuses[0].State.Waiting != nil {
					status = fmt.Sprintf("Waiting (%s)", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
				} else if pod.Status.ContainerStatuses[0].State.Running != nil {
					status = "Running"
				} else if pod.Status.ContainerStatuses[0].State.Terminated != nil {
					status = "Terminated"
				}

				healthy := "False"
				if pod.Status.ContainerStatuses[0].Ready {
					healthy = "True"
				}

				s := StatusOutput{
					Name:      label,
					Namespace: namespace,
					Created:   created,
					Age:       age,
					Status:    status,
					Version:   version,
					Healthy:   healthy,
				}

				m.Lock()
				statuses = append(statuses, s)
				m.Unlock()
			}
			wg.Done()
		}(lbl)
	}

	wg.Wait()
	return statuses
}
