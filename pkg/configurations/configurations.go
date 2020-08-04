package configurations

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/dapr/cli/pkg/standalone"
	v1alpha1 "github.com/dapr/dapr/pkg/apis/configuration/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"github.com/dapr/dashboard/pkg/age"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Configurations is an interface to interact with Dapr configurations
type Configurations interface {
	Supported() bool
	GetConfigurations() []Configuration
	GetConfiguration(name string) Configuration
}

type configurations struct {
	platform            string
	daprClient          scheme.Interface
	getConfigurationsFn func() []Configuration
}

// NewConfigurations returns a new Configurations instance
func NewConfigurations(platform string, daprClient scheme.Interface) Configurations {
	c := configurations{}
	c.platform = platform

	if platform == "kubernetes" {
		c.getConfigurationsFn = c.getKubernetesConfigurations
		c.daprClient = daprClient
	} else if platform == "standalone" {
		c.getConfigurationsFn = c.getStandaloneConfigurations
	}
	return &c
}

// Configuration represents a Dapr configuration
type Configuration struct {
	Name            string      `json:"name"`
	Kind            string      `json:"kind"`
	Created         string      `json:"created"`
	Age             string      `json:"age"`
	TracingEnabled  bool        `json:"tracingEnabled"`
	SamplingRate    string      `json:"samplingRate"`
	MTLSEnabled     bool        `json:"mtlsEnabled"`
	WorkloadCertTTL string      `json:"mtlsWorkloadTTL"`
	ClockSkew       string      `json:"mtlsClockSkew"`
	Manifest        interface{} `json:"manifest"`
}

// Supported checks whether or not the supplied platform is able to access Dapr configurations
func (c *configurations) Supported() bool {
	return c.platform == "kubernetes" || c.platform == "standalone"
}

// GetConfiguration returns the Dapr configuration specified by name
func (c *configurations) GetConfiguration(name string) Configuration {
	confs := c.getConfigurationsFn()
	for _, conf := range confs {
		if conf.Name == name {
			return conf
		}
	}
	return Configuration{}
}

// GetConfigurations returns the result of the correct platform's getConfigurations function
func (c *configurations) GetConfigurations() []Configuration {
	return c.getConfigurationsFn()
}

// getKubernetesConfigurations returns the list of all Dapr Configurations in a Kubernetes cluster
func (c *configurations) getKubernetesConfigurations() []Configuration {
	confs, err := c.daprClient.ConfigurationV1alpha1().Configurations(meta_v1.NamespaceAll).List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []Configuration{}
	}

	out := []Configuration{}
	for _, comp := range confs.Items {
		newConfiguration := Configuration{
			Name:            comp.Name,
			Kind:            comp.Kind,
			Created:         comp.CreationTimestamp.Format("2006-01-02 15:04.05"),
			Age:             age.GetAge(comp.CreationTimestamp.Time),
			TracingEnabled:  tracingEnabled(comp.Spec.TracingSpec),
			SamplingRate:    comp.Spec.TracingSpec.SamplingRate,
			MTLSEnabled:     comp.Spec.MTLSSpec.Enabled,
			WorkloadCertTTL: comp.Spec.MTLSSpec.WorkloadCertTTL,
			ClockSkew:       comp.Spec.MTLSSpec.AllowedClockSkew,
			Manifest:        comp,
		}
		out = append(out, newConfiguration)
	}
	return out
}

// getStandaloneConfigurations returns the list of Dapr Configurations Statuses
func (c *configurations) getStandaloneConfigurations() []Configuration {
	configurationFile := standalone.DefaultConfigFilePath()
	content, err := ioutil.ReadFile(configurationFile)
	if err != nil {
		log.Printf("Failure reading file %s: %v\n", configurationFile, err)
		return []Configuration{}
	}

	info, err := os.Stat(configurationFile)
	if err != nil {
		log.Printf("Failure reading file info from %s: %v\n", configurationFile, err)
		return []Configuration{}
	}

	comp := v1alpha1.Configuration{}
	err = yaml.Unmarshal(content, &comp)
	if err != nil {
		fmt.Println(err.Error())
	}

	newConfiguration := Configuration{
		Name:            comp.Name,
		Kind:            comp.Kind,
		Created:         info.ModTime().Format("2006-01-02 15:04.05"),
		Age:             age.GetAge(info.ModTime()),
		TracingEnabled:  tracingEnabled(comp.Spec.TracingSpec),
		SamplingRate:    comp.Spec.TracingSpec.SamplingRate,
		MTLSEnabled:     comp.Spec.MTLSSpec.Enabled,
		WorkloadCertTTL: comp.Spec.MTLSSpec.WorkloadCertTTL,
		ClockSkew:       comp.Spec.MTLSSpec.AllowedClockSkew,
		Manifest:        string(content),
	}

	standaloneConfigurations := []Configuration{}
	return append(standaloneConfigurations, newConfiguration)
}

// tracingEnabled checks if tracing is enabled for a configuration
func tracingEnabled(spec v1alpha1.TracingSpec) bool {
	sr, err := strconv.ParseFloat(spec.SamplingRate, 32)
	if err != nil {
		return false
	}
	return sr > 0
}
