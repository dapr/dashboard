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
	address := flag.String("address", "localhost", "Address to listen on. Only accepts IP address or localhost as a value")
	port := flag.Int("port", 8080, "Port to listen to")

	flag.Parse()

	if *dashboardVersion {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	} else {
		RunWebServer(*address, *port)
	}
}
