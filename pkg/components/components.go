package components

import (
	"log"

	v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"github.com/dapr/dashboard/pkg/age"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Components is an interface to interact with Dapr components
type Components interface {
	Supported() bool
	GetComponents(scope string) []v1alpha1.Component
	GetComponent(scope string, name string) v1alpha1.Component
	GetStatus(scope string) []ComponentsOutput
}

type components struct {
	daprClient scheme.Interface
}

// ComponentsOutput represent a Dapr component.
type ComponentsOutput struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Age     string `json:"age"`
	Created string `json:"created"`
}

// NewComponents returns a new Components instance
func NewComponents(daprClient scheme.Interface) Components {
	return &components{
		daprClient: daprClient,
	}
}

// Supported checks whether or not the dapr client is available to access the components
func (c *components) Supported() bool {
	return c.daprClient != nil
}

// GetComponents returns the list of all Dapr components
func (c *components) GetComponents(scope string) []v1alpha1.Component {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(scope).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []v1alpha1.Component{}
	}
	return comps.Items
}

// GetComponent returns a specific component based on a supplied component name
func (c *components) GetComponent(scope string, name string) v1alpha1.Component {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(scope).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return v1alpha1.Component{}
	}
	for _, c := range comps.Items {
		if c.ObjectMeta.Name == name {
			return c
		}
	}
	return v1alpha1.Component{}
}

// GetStatus returns returns a list of Dapr component statuses
func (c *components) GetStatus(scope string) []ComponentsOutput {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(scope).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []ComponentsOutput{}
	}

	co := []ComponentsOutput{}
	for _, c := range comps.Items {
		co = append(co, ComponentsOutput{
			Name:    c.GetName(),
			Type:    c.Spec.Type,
			Created: c.CreationTimestamp.Format("2006-01-02 15:04.05"),
			Age:     age.GetAge(c.CreationTimestamp.Time),
		})
	}
	return co
}
