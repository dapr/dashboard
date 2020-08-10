package instances

// Instance describes a Dapr sidecar instance information
type Instance struct {
	AppID            string `json:"appID"`
	HTTPPort         int    `json:"httpPort"`
	GRPCPort         int    `json:"grpcPort"`
	AppPort          int    `json:"appPort"`
	Command          string `json:"command"`
	Age              string `json:"age"`
	Created          string `json:"created"`
	PID              int    `json:"pid"`
	Replicas         int    `json:"replicas"`
	Address          string `json:"address"`
	SupportsDeletion bool   `json:"supportsDeletion"`
	SupportsLogs     bool   `json:"supportsLogs"`
	Manifest         string `json:"manifest"`
	Status           string `json:"status"`
	Labels           string `json:"labels"`
	Selector         string `json:"selector"`
	Config           string `json:"config"`
}

// StatusOutput represents the status of a named Dapr resource.
type StatusOutput struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Healthy   string `json:"healthy"`
	Status    string `json:"status"`
	Version   string `json:"version"`
	Age       string `json:"age"`
	Created   string `json:"created"`
}

// MetadataOutput represents a metadata api call response
type MetadataOutput struct {
	ID       string                      `json:"id"`
	Actors   []MetadataActiveActorsCount `json:"actors"`
	Extended map[string]interface{}      `json:"extended"`
}

// MetadataActiveActorsCount represents actor metadata: type and count
type MetadataActiveActorsCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// Log represents a log message with metadata
type Log struct {
	Level     string `json:"level"`
	Timestamp int64  `json:"timestamp"`
	Container string `json:"container"`
	Content   string `json:"content"`
}
