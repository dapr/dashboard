package configurations

import (
	"log"
	"strconv"

	v1alpha1 "github.com/dapr/dapr/pkg/apis/configuration/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"github.com/dapr/dashboard/pkg/age"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Configurations is an interface to interact with a Dapr configuration
type Configurations interface {
	Supported() bool
	GetConfigurations(scope string) []v1alpha1.Configuration
	GetConfiguration(name string, scope string) v1alpha1.Configuration
	GetStatus(scope string) []ConfigurationsOutput
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

// Supported checks whether or not the dapr client is available to access the configurations
func (c *configurations) Supported() bool {
	return c.daprClient != nil
}

// GetConfigurations returns the list of all Dapr Configurations
func (c *configurations) GetConfigurations(scope string) []v1alpha1.Configuration {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(scope).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []v1alpha1.Configuration{}
	}

	return confs.Items
}

// GetConfiguration returns one Dapr configuration
func (c *configurations) GetConfiguration(scope string, name string) v1alpha1.Configuration {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(scope).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return v1alpha1.Configuration{}
	}

	for _, c := range confs.Items {
		if c.ObjectMeta.Name == name {
			return c
		}
	}
	return v1alpha1.Configuration{}
}

// GetStatus returns the list of Dapr Configurations Statuses
func (c *configurations) GetStatus(scope string) []ConfigurationsOutput {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(scope).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []ConfigurationsOutput{}
	}

	co := []ConfigurationsOutput{}
	for _, c := range confs.Items {
		co = append(co, ConfigurationsOutput{
			TracingEnabled:  tracingEnabled(c.Spec.TracingSpec),
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

// tracingEnabled checks if tracing is enabled for a configuration
func tracingEnabled(spec v1alpha1.TracingSpec) bool {
	sr, err := strconv.ParseFloat(spec.SamplingRate, 32)
	if err != nil {
		return false
	}
	return sr > 0
}
