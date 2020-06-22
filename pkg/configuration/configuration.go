package configuration

import (
	"github.com/dapr/cli/pkg/age"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Configuration is an interface to interact with a Dapr configuration
type Configuration interface {
	Supported() bool
	Get() []ConfigurationsOutput
}

type configuration struct {
	daprClient scheme.Interface
}

// ConfigurationsOutput represent a Dapr configuration
type ConfigurationsOutput struct {
	Name            string `csv:"Name"`
	TracingEnabled  bool   `csv:"TRACING-ENABLED"`
	MTLSEnabled     bool   `csv:"MTLS-ENABLED"`
	WorkloadCertTTL string `csv:"MTLS-WORKLOAD-TTL"`
	ClockSkew       string `csv:"MTLS-CLOCK-SKEW"`
	Age             string `csv:"AGE"`
	Created         string `csv:"CREATED"`
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

func (c *configuration) Get() []ConfigurationsOutput {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(meta_v1.NamespaceAll).List(meta_v1.ListOptions{})
	if err != nil {
		return []ConfigurationsOutput{}
	}

	co := []ConfigurationsOutput{}
	for _, c := range confs.Items {
		co = append(co, ConfigurationsOutput{
			TracingEnabled:  c.Spec.TracingSpec.Enabled,
			Name:            c.GetName(),
			MTLSEnabled:     c.Spec.MTLSSpec.Enabled,
			WorkloadCertTTL: c.Spec.MTLSSpec.WorkloadCertTTL,
			ClockSkew:       c.Spec.MTLSSpec.AllowedClockSkew,
			Created:         c.CreationTimestamp.Format("2006-01-02 15:04.05"),
			Age:             age.GetAge(c.CreationTimestamp.Time),
		})
	}
	return co
}
