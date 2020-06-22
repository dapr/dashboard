package configurations

import (
	"log"

	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"github.com/dapr/dashboard/pkg/age"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Configurations is an interface to interact with a Dapr configuration
type Configurations interface {
	Supported() bool
	Get() []ConfigurationsOutput
}

type configurations struct {
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

// NewConfigurations returns a new Configurations instance
func NewConfigurations(daprClient scheme.Interface) Configurations {
	return &configurations{
		daprClient: daprClient,
	}
}

func (c *configurations) Supported() bool {
	return c.daprClient != nil
}

func (c *configurations) Get() []ConfigurationsOutput {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(meta_v1.NamespaceAll).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []ConfigurationsOutput{}
	}

	co := []ConfigurationsOutput{}
	for _, c := range confs.Items {
		configs = append(configs, ConfigurationsOutput{
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
