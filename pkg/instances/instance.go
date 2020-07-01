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
}
