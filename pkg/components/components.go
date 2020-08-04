package components

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dapr/cli/pkg/standalone"
	v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"github.com/dapr/dashboard/pkg/age"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Components is an interface to interact with Dapr components
type Components interface {
	Supported() bool
	GetComponents() []Component
	GetComponent(name string) Component
}

type components struct {
	platform        string
	daprClient      scheme.Interface
	getComponentsFn func() []Component
}

// NewComponents returns a new Components instance
func NewComponents(platform string, daprClient scheme.Interface) Components {
	c := components{}
	c.platform = platform

	if platform == "kubernetes" {
		c.getComponentsFn = c.getKubernetesComponents
		c.daprClient = daprClient
	} else if platform == "standalone" {
		c.getComponentsFn = c.getStandaloneComponents
	}
	return &c
}

// Component represents a platform-independent Dapr component
type Component struct {
	Name     string      `json:"name"`
	Kind     string      `json:"kind"`
	Type     string      `json:"type"`
	Created  string      `json:"created"`
	Age      string      `json:"age"`
	Scopes   []string    `json:"scopes"`
	Manifest interface{} `json:"manifest"`
}

// Supported checks whether or not the supplied platform is able to access Dapr components
func (c *components) Supported() bool {
	return c.platform == "kubernetes" || c.platform == "standalone"
}

// GetComponent returns a specific component based on a supplied component name
func (c *components) GetComponent(name string) Component {
	comps := c.getComponentsFn()
	for _, comp := range comps {
		if comp.Name == name {
			return comp
		}
	}
	return Component{}
}

// GetComponent returns the result of the correct platform's getComponents function
func (c *components) GetComponents() []Component {
	return c.getComponentsFn()
}

// getKubernetesComponents returns the list of all Dapr components in a Kubernetes cluster
func (c *components) getKubernetesComponents() []Component {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(metav1.NamespaceDefault).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []Component{}
	}
	out := []Component{}
	for _, comp := range comps.Items {
		newComponent := Component{
			Name:     comp.Name,
			Kind:     comp.Kind,
			Type:     comp.Spec.Type,
			Created:  comp.CreationTimestamp.Format("2006-01-02 15:04.05"),
			Age:      age.GetAge(comp.CreationTimestamp.Time),
			Scopes:   comp.Scopes,
			Manifest: comp,
		}
		out = append(out, newComponent)
	}
	return out
}

func (c *components) getStandaloneComponents() []Component {
	componentsDirectory := standalone.DefaultComponentsDirPath()
	standaloneComponents := []Component{}
	err := filepath.Walk(componentsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failure accessing path %s: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Failure reading file %s: %v\n", path, err)
				return err
			}

			comp := v1alpha1.Component{}
			err = yaml.Unmarshal(content, &comp)
			if err != nil {
				fmt.Println(err.Error())
			}

			newComponent := Component{
				Name:     comp.Name,
				Kind:     comp.Kind,
				Type:     comp.Spec.Type,
				Created:  info.ModTime().Format("2006-01-02 15:04.05"),
				Age:      age.GetAge(info.ModTime()),
				Scopes:   comp.Scopes,
				Manifest: string(content),
			}
			standaloneComponents = append(standaloneComponents, newComponent)
			return nil
		}
		return nil
	})
	if err != nil {
		log.Printf("error walking the path %q: %v\n", componentsDirectory, err)
		return []Component{}
	}
	return standaloneComponents
}
