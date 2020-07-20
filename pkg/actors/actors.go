package actors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Actors is an interface to interact with Dapr actors
type Actors interface {
	Get(id string, port int) string
}

type actors struct {
	kubeClient *kubernetes.Clientset
}

const (
	daprEnabledAnnotation = "dapr.io/enabled"
	daprIDAnnotation      = "dapr.io/id"
	daprPortAnnotation    = "dapr.io/port"
)

// ActorsOutput represents an Actor api call response
type ActorsOutput struct {
	ID       string                      `json:"id`
	Actors   []MetadataActiveActorsCount `json:"actors"`
	Extended map[string]interface{}      `json:"extended"`
}

// MetadataActiveActorsCount represents actor metadata: type and count
type MetadataActiveActorsCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// NewActors returns a new Actors instance
func NewActors(kubeClient *kubernetes.Clientset) Actors {
	return &actors{
		kubeClient: kubeClient,
	}
}

func (a *actors) Get(id string, port int) string {
	if a.kubeClient != nil {
		resp, err := a.kubeClient.AppsV1().Deployments(meta_v1.NamespaceAll).List((meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return ""
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := a.kubeClient.CoreV1().Pods(d.GetNamespace()).List(meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						log.Println(err)
						return ""
					}

					if len(pods.Items) > 0 {
						p := pods.Items[0]

						url := fmt.Sprintf("http://%v:%v/v1.0/metadata", p.Status.PodIP, 3500)

						resp, err := http.Get(url)
						if err != nil {
							log.Println(err)
							return ""
						}

						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							log.Println(err)
							return ""
						}

						var data ActorsOutput

						if err := json.Unmarshal(body, &data); err != nil {
							return string(body)
						}

						out := ""
						for _, act := range data.Actors {
							out += fmt.Sprintf("{id: %s, count: %v}\n", act.Type, act.Count)
						}

						return out
					}
				}
			}
		}

	} else {
		url := fmt.Sprintf("http://localhost:%v/v1.0/metadata", port)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return ""
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return ""
		}

		var data ActorsOutput

		if err := json.Unmarshal(body, &data); err != nil {
			return string(body)
		}

		out := ""
		for _, act := range data.Actors {
			actorJSON, err := json.Marshal(act)
			if err != nil {
				log.Println(err)
				return ""
			}
			out += string(actorJSON)
		}

		return out
	}
	return ""
}
