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
