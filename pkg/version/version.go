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

package version

import (
	"github.com/dapr/cli/pkg/standalone"
	"strings"
)

// version is the current Dapr dashboard version
var version = "edge"

// GetVersion returns the current dashboard version
func GetVersion() string {
	return version
}

// GetRuntimeVersion returns the current runtime version
func GetRuntimeVersion() string {
	return strings.ReplaceAll(standalone.GetRuntimeVersion(), "\n", "")
}
