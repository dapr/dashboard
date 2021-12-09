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
	Service   string `json:"service"`
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
