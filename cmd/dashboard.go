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

	flag.Parse()

	if *dashboardVersion {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	} else {
		RunWebServer()
	}
}
