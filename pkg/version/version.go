package version

// version is the current Dapr dashboard version
var version = "edge"

// GetVersion returns the current dashboard version
func GetVersion() string {
	return version
}
