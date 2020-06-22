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
	Get() []v1alpha1.Component
	GetStatus() []ComponentsOutput
}

type components struct {
	daprClient scheme.Interface
}

// ComponentsOutput represent a Dapr component.
type ComponentsOutput struct {
	Name    string `csv:"Name"`
	Type    string `csv:"Type"`
	Age     string `csv:"AGE"`
	Created string `csv:"CREATED"`
}

// NewComponents returns a new Components instance
func NewComponents(daprClient scheme.Interface) Components {
	return &components{
		daprClient: daprClient,
	}
}

func (c *components) Supported() bool {
	return c.daprClient != nil
}

func (c *components) Get() []v1alpha1.Component {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(meta_v1.NamespaceDefault).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []v1alpha1.Component{}
	}
	return comps.Items
}

func (c *components) GetStatus() []ComponentsOutput {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(meta_v1.NamespaceDefault).List(meta_v1.ListOptions{})
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
