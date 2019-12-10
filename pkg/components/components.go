package components

import (
	"log"

	v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Components is an interface to interact with Dapr components
type Components interface {
	Supported() bool
	Get() []v1alpha1.Component
}

type components struct {
	daprClient scheme.Interface
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
