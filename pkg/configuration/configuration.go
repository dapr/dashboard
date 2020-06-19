package configuration

import (
	"github.com/dapr/cli/pkg/kubernetes"
	v1alpha1 "github.com/dapr/dapr/pkg/apis/configuration/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Configuration is an interface to interact with a Dapr configuration
type Configuration interface {
	Supported() bool
	Get() []v1alpha1.Configuration
}

type configuration struct {
	daprClient scheme.Interface
}

// NewConfiguration returns a new Configuration instance
func NewConfiguration(daprClient scheme.Interface) Configuration {
	return &configuration{
		daprClient: daprClient,
	}
}

func (c *configuration) Supported() bool {
	return c.daprClient != nil
}

func (c *configuration) Get() []v1alpha1.Configuration {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(meta_v1.NamespaceDefault).List(meta_v1.ListOptions{})
	if err != nil {
		// For debugging
		basic := []v1alpha1.Configuration{}
		basic = append(basic, kubernetes.GetDefaultConfiguration())
		basic[0].ObjectMeta.Name = err.Error()
		return basic
	}
	return confs.Items
}
