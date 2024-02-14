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

package platforms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlatforms(t *testing.T) {
	var scenarios = []struct {
		platform Platform
		want     string
	}{
		{Kubernetes, "kubernetes"},
		{Standalone, "standalone"},
		{DockerCompose, "docker-compose"},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Platform %s as string", scenario.platform), func(t *testing.T) {
			platformAsString := string(scenario.platform)
			assert.Equal(t, scenario.want, platformAsString)
		})
	}
}
