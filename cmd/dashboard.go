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

package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dapr/dashboard/pkg/version"
)

// RunDashboard runs the dashboard with the supplied flags
func RunDashboard() {
	dashboardVersion := flag.Bool("version", false, "Prints the dashboard version")
	port := flag.Int("port", 8080, "Port to listen to")
	dockerCompose := flag.Bool("docker-compose", false, "Is running inside docker compose")
	componentsPath := flag.String("components-path", "", "Path to volume mounted dapr components (docker-compose only)")
	configPath := flag.String("config-path", "", "Path to volume mounted dapr configuration (docker-compose only)")
	dockerComposePath := flag.String("docker-compose-path", "", "Path to volume mounted docker compose file (docker-compose only)")

	flag.Parse()

	fmt.Println(*dockerCompose)
	fmt.Println(*componentsPath)
	fmt.Println(*configPath)
	fmt.Println(*dockerComposePath)

	if *dockerCompose {
		if len(*componentsPath) == 0 {
			fmt.Println("--components-path required when --docker-compose=true")
			os.Exit(1)
		}

		if len(*configPath) == 0 {
			fmt.Println("--config-path required when --docker-compose=true")
			os.Exit(1)
		}

		if len(*dockerComposePath) == 0 {
			fmt.Println("--docker-compose-path required when --docker-compose=true")
			os.Exit(1)
		}
	}

	if *dashboardVersion {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	} else {
		RunWebServer(*port, *dockerCompose, *componentsPath, *configPath, *dockerComposePath)
	}
}
