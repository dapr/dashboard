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
package configurations

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
			target := NewConfigurations(scenario.platform, nil, "")
			isSupported := target.Supported()
			assert.Equal(t, scenario.want, isSupported)
		})
	}
}

func TestDockerComposeGetConfigurations(t *testing.T) {
	target := NewConfigurations(platforms.DockerCompose, nil, "testdata")
	configurations := target.GetConfigurations("")
	assert.Equal(t, 1, len(configurations), "Should load test data configurations")
	for _, v := range configurations {
		assert.Equal(t, v.Kind, "Configuration", "Should only load configuration files")
	}
}

func TestDockerComposeGetConfiguration(t *testing.T) {
	var scenarios = []struct {
		name string
		want bool
	}{
		{"tracing", true},
		{"does_not_exist", false},
	}

	target := NewConfigurations(platforms.DockerCompose, nil, "testdata")

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Should load valid configuration - %t", scenario.want), func(t *testing.T) {
			config := target.GetConfiguration("", scenario.name)
			assert.NotNil(t, config, "Should always return something")

			if scenario.want {
				assert.Equal(t, "Configuration", config.Kind, "Should only return configurations")
				assert.Equal(t, scenario.name, config.Name, "When configuration valid, name is set")
			} else {
				assert.Empty(t, config.Kind, "When configuration not valid, kind is not set")
				assert.Empty(t, config.Name, "When configuration not valid, name is not set")
			}
		})
	}
}
