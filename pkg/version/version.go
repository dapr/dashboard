package version

// Version is the current Dapr dashboard version
var Version = "edge"

// GetVersion returns the current dashboard version
func GetVersion() string {
	return Version
}
