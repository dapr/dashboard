package actors

import (
	"bytes"
	"log"

	scheme "github.com/dapr/dapr/pkg/client/clientset/versioned"
)

// Actors is an interface to interact with Dapr actors
type Actors interface {
	Supported() bool
	Get() string
}

type actors struct {
	daprClient scheme.Interface
}

// MetadataActiveActorsCount represents actor metadata: type and count
type MetadataActiveActorsCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// NewActors returns a new Actors instance
func NewActors(daprClient scheme.Interface) Actors {
	return &actors{
		daprClient: daprClient,
	}
}

// Supported checks whether or not the Dapr client is available
func (a *actors) Supported() bool {
	return a.daprClient != nil
}

func (a *actors) Get() string {
	if a.Supported() {
		b := a.daprClient
		c := b.ComponentsV1alpha1()
		restClient := c.RESTClient()
		// restClient := a.daprClient.Discovery().RESTClient()
		data, err := restClient.Get().RequestURI("/v1.0/metadata").Stream()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(data)
		if err != nil {
			log.Println(err)
			return ""
		}
		dataStr := buf.String()

		return dataStr
	}
	return ""
}
