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
	"fmt"
	"testing"

	"github.com/dapr/dashboard/pkg/platforms"
	"github.com/stretchr/testify/assert"
)

func TestSupported(t *testing.T) {
	var scenarios = []struct {
		platform platforms.Platform
		want     bool
	}{
		{platforms.Kubernetes, true},
		{platforms.Standalone, true},
		{platforms.DockerCompose, true},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Platform %s should be supported", scenario.platform), func(t *testing.T) {
			target := NewComponents(scenario.platform, nil, "")
			isSupported := target.Supported()
			assert.Equal(t, scenario.want, isSupported)
		})
	}
}

func TestDockerComposeGetComponents(t *testing.T) {
	target := NewComponents(platforms.DockerCompose, nil, "testdata")
	components := target.GetComponents("")
	assert.Equal(t, 2, len(components), "Should load test data components")
	for _, v := range components {
		assert.Equal(t, v.Kind, "Component", "Should only load component files")
	}
}

func TestDockerComposeGetComponent(t *testing.T) {
	var scenarios = []struct {
		name string
		want bool
	}{
		{"cronjob", true},
		{"messagebus", true},
		{"does_not_exist", false},
	}

	target := NewComponents(platforms.DockerCompose, nil, "testdata")

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Should load valid component - %t", scenario.want), func(t *testing.T) {
			component := target.GetComponent("", scenario.name)
			assert.NotNil(t, component, "Should always return something")

			if scenario.want {
				assert.Equal(t, "Component", component.Kind, "Should only return components")
				assert.Equal(t, scenario.name, component.Name, "Name should be set")
				assert.NotEmpty(t, component.Type, "When component valid, type is set")
			} else {
				assert.Empty(t, component.Kind, "When component not valid, kind is not set")
				assert.Empty(t, component.Type, "When component not valid, type is not set")
			}
		})
	}
}

func TestIsYamlFile(t *testing.T) {
	var scenarios = []struct {
		path string
		want bool
	}{
		{"pubsub.yaml", true},
		{"pubsub.yml", true},
		{"pubsub.txt", false},
		{"pubsub", false},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Should return valid yaml file - %t", scenario.want), func(t *testing.T) {
			actual := isYamlFile(scenario.path)
			assert.Equal(t, scenario.want, actual)
		})
	}
}
