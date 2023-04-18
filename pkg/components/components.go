/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package components

import (
	"log"
	"os"
	"path/filepath"

	"github.com/dapr/cli/pkg/standalone"
	v1alpha1 "github.com/dapr/dapr/pkg/apis/components/v1alpha1"
	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
	"github.com/dapr/dashboard/pkg/age"
	"github.com/dapr/dashboard/pkg/platforms"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Components is an interface to interact with Dapr components
type Components interface {
	Supported() bool
	GetComponent(scope string, name string) Component
	GetComponents(scope string) []Component
}

type components struct {
	platform        platforms.Platform
	daprClient      scheme.Interface
	getComponentsFn func(scope string) []Component
	componentsPath  string
}

// NewComponents returns a new Components instance
func NewComponents(platform platforms.Platform, daprClient scheme.Interface, componentsPath string) Components {
	c := components{}
	c.platform = platform

	if platform == platforms.Kubernetes {
		c.getComponentsFn = c.getKubernetesComponents
		c.daprClient = daprClient
	} else if platform == platforms.Standalone {
		c.getComponentsFn = c.getStandaloneComponents
	} else if platform == platforms.DockerCompose {
		c.getComponentsFn = c.getDockerComposeComponents
		c.componentsPath = componentsPath
	}
	return &c
}

// Component represents a Dapr component
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
	return c.platform == platforms.Kubernetes || c.platform == platforms.Standalone || c.platform == platforms.DockerCompose
}

// GetComponent returns a specific component based on a supplied component name
func (c *components) GetComponent(scope string, name string) Component {
	comps := c.getComponentsFn(scope)
	for _, comp := range comps {
		if comp.Name == name {
			return comp
		}
	}
	return Component{}
}

// GetComponent returns the result of the correct platform's getComponents function
func (c *components) GetComponents(scope string) []Component {
	return c.getComponentsFn(scope)
}

// getKubernetesComponents returns the list of all Dapr components in a Kubernetes cluster
func (c *components) getKubernetesComponents(scope string) []Component {
	comps, err := c.daprClient.ComponentsV1alpha1().Components(scope).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return []Component{}
	}
	out := []Component{}
	for _, comp := range comps.Items {
		newComponent := Component{
			Name:     comp.Name,
			Kind:     comp.Kind(),
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

// getStandaloneComponents returns the list of all locally-hosted Dapr components
func (c *components) getStandaloneComponents(scope string) []Component {
	standaloneComponents := []Component{}

	daprDir, err := standalone.GetDaprRuntimePath("")
	if err != nil {
		log.Printf("Failure findinf Dapr's runtime path: %v\n", err)
		return standaloneComponents
	}

	componentsDirectory := standalone.GetDaprComponentsPath(daprDir)

	err = filepath.Walk(componentsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failure accessing path %s: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() != filepath.Base(componentsDirectory) {
			return filepath.SkipDir
		} else if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failure reading file %s: %v\n", path, err)
				return err
			}

			comp := v1alpha1.Component{}
			err = yaml.Unmarshal(content, &comp)
			if err != nil {
				log.Printf("Failure unmarshalling %s into Component: %s\n", path, err.Error())
			}

			newComponent := Component{
				Name:     comp.Name,
				Kind:     comp.Kind(),
				Type:     comp.Spec.Type,
				Created:  info.ModTime().Format("2006-01-02 15:04.05"),
				Age:      age.GetAge(info.ModTime()),
				Scopes:   comp.Scopes,
				Manifest: string(content),
			}

			if newComponent.Kind == "Component" {
				standaloneComponents = append(standaloneComponents, newComponent)
			}
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

// getDockerComposeComponents returns the list of all docker-compose Dapr components
func (c *components) getDockerComposeComponents(scope string) []Component {
	componentsDirectory := c.componentsPath
	dockerComposeComponents := []Component{}
	err := filepath.Walk(componentsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failure accessing path %s: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() != filepath.Base(componentsDirectory) {
			return filepath.SkipDir
		} else if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Failure reading file %s: %v\n", path, err)
				return err
			}

			comp := v1alpha1.Component{}
			err = yaml.Unmarshal(content, &comp)
			if err != nil {
				log.Printf("Failure unmarshalling %s into Component: %s\n", path, err.Error())
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

			if newComponent.Kind == "Component" {
				dockerComposeComponents = append(dockerComposeComponents, newComponent)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		log.Printf("error walking the path %q: %v\n", componentsDirectory, err)
		return []Component{}
	}
	return dockerComposeComponents
}
