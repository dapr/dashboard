package instances

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dapr/cli/pkg/standalone"
	"github.com/dapr/dashboard/pkg/age"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	json_serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

const (
	daprEnabledAnnotation = "dapr.io/enabled"
	daprIDAnnotation      = "dapr.io/id"
	daprPortAnnotation    = "dapr.io/port"
)

var (
	controlPlaneLabels = [...]string{"dapr-operator", "dapr-sentry", "dapr-placement", "dapr-sidecar-injector"}
)

// Instances is an interface to interact with running Dapr instances in Kubernetes or Standalone modes
type Instances interface {
	Supported() bool
	GetInstances(scope string) []Instance
	GetInstance(scope string, id string) Instance
	DeleteInstance(scope string, id string) error
	GetLogs(scope string, id string) []Log
	GetDeploymentConfiguration(scope string, id string) string
	GetControlPlaneStatus() []StatusOutput
	GetMetadata(scope string, id string) MetadataOutput
	GetActiveActorsCount(metadata MetadataOutput) []MetadataActiveActorsCount
	GetScopes() []string
	CheckPlatform() string
}

type instances struct {
	platform       string
	kubeClient     *kubernetes.Clientset
	getInstancesFn func(string) []Instance
	getScopesFn    func() []string
}

// NewInstances returns an Instances instance
func NewInstances(platform string, kubeClient *kubernetes.Clientset) Instances {
	i := instances{}
	i.platform = platform

	if i.platform == "kubernetes" {
		i.getInstancesFn = i.getKubernetesInstances
		i.getScopesFn = i.getKubernetesScopes
		i.kubeClient = kubeClient
	} else if i.platform == "standalone" {
		i.getInstancesFn = i.getStandaloneInstances
		i.getScopesFn = i.getStandaloneScopes
	}
	return &i
}

// Supported checks if the current platform supports Dapr instances
func (i *instances) Supported() bool {
	return i.platform == "kubernetes" || i.platform == "standalone"
}

// GetScopes returns the result of the appropriate environment's GetScopes function
func (i *instances) GetScopes() []string {
	return i.getScopesFn()
}

// CheckPlatform returns the current environment dashboard is running in
func (i *instances) CheckPlatform() string {
	return i.platform
}

// GetLogs returns a string of all logs for the given Dapr app id
func (i *instances) GetLogs(scope string, id string) []Log {
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List((meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return []Log{}
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						log.Println(err)
						return []Log{}
					}

					if len(pods.Items) > 0 {
						p := pods.Items[0]
						name := p.ObjectMeta.Name

						out := []Log{}
						for _, container := range p.Spec.Containers {
							options := v1.PodLogOptions{}
							options.Container = container.Name
							options.Timestamps = true

							res := i.kubeClient.CoreV1().Pods(p.ObjectMeta.Namespace).GetLogs(name, &options)
							stream, err := res.Stream()
							if err != nil {
								log.Println(err)
								return out
							}

							buf := new(bytes.Buffer)
							_, err = buf.ReadFrom(stream)
							if err != nil {
								log.Println(err)
								return out
							}
							bufString := buf.String()

							levelExp, _ := regexp.Compile("(level=)[^ ]*")
							timeExp, _ := regexp.Compile("^[^ ]+")

							for _, content := range strings.Split(bufString, "\n") {
								currentLog := Log{
									Level:     "info",
									Timestamp: 0,
									Container: container.Name,
									Content:   content,
								}

								currentLog.Level = strings.Replace(levelExp.FindString(content), "level=", "", 1)
								if err != nil {
									log.Println(err)
									continue
								}

								timestamp, err := time.Parse(time.RFC3339Nano, timeExp.FindString(content))
								if err != nil {
									log.Println(err)
									continue
								}

								currentLog.Timestamp = timestamp.UnixNano()
								out = append(out, currentLog)
							}
						}
						return out
					}
				}
			}
		}
	}
	return []Log{}
}

// GetDeploymentConfiguration returns the metadata of a Dapr application in YAML format
func (i *instances) GetDeploymentConfiguration(scope string, id string) string {
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List((meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return ""
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						log.Println(err)
						return ""
					}

					if len(pods.Items) > 0 {
						p := pods.Items[0]

						name := p.ObjectMeta.Name
						nspace := p.ObjectMeta.Namespace

						restClient := i.kubeClient.CoreV1().RESTClient()
						if err != nil {
							log.Println(err)
							return ""
						}

						url := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", nspace, name)
						data, err := restClient.Get().RequestURI(url).Stream()
						if err != nil {
							log.Println(err)
							return ""
						}

						buf := new(bytes.Buffer)
						_, err = buf.ReadFrom(data)
						if err != nil {
							log.Println(err)
							return ""
						}
						dataStr := buf.String()
						j := []byte(dataStr)
						y, err := yaml.JSONToYAML(j)
						if err != nil {
							log.Println(err)
							return ""
						}

						return string(y)
					}
				}
			}
		}

	}
	return ""
}

// DeleteInstance deletes the local Dapr sidecar instance
func (i *instances) DeleteInstance(scope string, id string) error {
	return standalone.Stop(id)
}

// GetInstance uses the appropriate getInstance function (kubernetes, standalone, etc.) and returns the given instance from its id
func (i *instances) GetInstance(scope string, id string) Instance {
	instanceList := i.getInstancesFn(scope)
	for _, instance := range instanceList {
		if instance.AppID == id {
			return instance
		}
	}
	return Instance{}
}

// GetControlPlaneStatus lists the status of each of the Dapr control plane services
func (i *instances) GetControlPlaneStatus() []StatusOutput {
	if i.kubeClient != nil {
		var wg sync.WaitGroup
		wg.Add(len(controlPlaneLabels))

		m := sync.Mutex{}
		statuses := []StatusOutput{}

		for _, lbl := range controlPlaneLabels {
			go func(label string) {
				options := meta_v1.ListOptions{}
				labelSelector := map[string]string{
					"app": label,
				}
				options.LabelSelector = labels.FormatLabels(labelSelector)

				p, err := i.kubeClient.CoreV1().Pods(v1.NamespaceAll).List(options)
				if err == nil && len(p.Items) == 1 {
					pod := p.Items[0]

					image := pod.Spec.Containers[0].Image
					namespace := pod.GetNamespace()
					age := age.GetAge(pod.CreationTimestamp.Time)
					created := pod.CreationTimestamp.Format("2006-01-02 15:04.05")
					version := image[strings.IndexAny(image, ":")+1:]
					status := ""

					if pod.Status.ContainerStatuses[0].State.Waiting != nil {
						status = fmt.Sprintf("Waiting (%s)", pod.Status.ContainerStatuses[0].State.Waiting.Reason)
					} else if pod.Status.ContainerStatuses[0].State.Running != nil {
						status = "Running"
					} else if pod.Status.ContainerStatuses[0].State.Terminated != nil {
						status = "Terminated"
					}

					healthy := "False"
					if pod.Status.ContainerStatuses[0].Ready {
						healthy = "True"
					}

					s := StatusOutput{
						Name:      label,
						Namespace: namespace,
						Created:   created,
						Age:       age,
						Status:    status,
						Version:   version,
						Healthy:   healthy,
					}

					m.Lock()
					statuses = append(statuses, s)
					m.Unlock()
				}
				wg.Done()
			}(lbl)
		}

		wg.Wait()
		return statuses
	}
	return []StatusOutput{}
}

// GetMetadata returns the result from the /v1.0/metadata endpoint
func (i *instances) GetMetadata(scope string, id string) MetadataOutput {
	url := ""
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List((meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return MetadataOutput{}
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						log.Println(err)
						return MetadataOutput{}
					}

					if len(pods.Items) > 0 {
						p := pods.Items[0]
						url = fmt.Sprintf("http://%v:%v/v1.0/metadata", p.Status.PodIP, 3500)
					}
				}
			}
		}

	} else {
		port := i.GetInstance(scope, id).HTTPPort
		url = fmt.Sprintf("http://localhost:%v/v1.0/metadata", port)
	}
	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return MetadataOutput{}
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return MetadataOutput{}
		}

		var data MetadataOutput
		if err := json.Unmarshal(body, &data); err != nil {
			log.Println(err)
			return MetadataOutput{}
		}
		return data
	}
	return MetadataOutput{}
}

// GetActiveActorsCount returns the Actors slice of a MetadataOutput
func (i *instances) GetActiveActorsCount(metadata MetadataOutput) []MetadataActiveActorsCount {
	return metadata.Actors
}

// GetInstances returns the result of the appropriate environment's GetInstance function
func (i *instances) GetInstances(scope string) []Instance {
	return i.getInstancesFn(scope)
}

// getKubernetesInstances gets the list of Dapr applications running in the Kubernetes environment
func (i *instances) getKubernetesInstances(scope string) []Instance {
	list := []Instance{}
	resp, err := i.kubeClient.AppsV1().Deployments(scope).List((meta_v1.ListOptions{}))
	if err != nil {
		log.Println(err)
		return list
	}

	for _, d := range resp.Items {
		if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
			id := d.Spec.Template.Annotations[daprIDAnnotation]
			i := Instance{
				AppID:            id,
				HTTPPort:         3500,
				GRPCPort:         50001,
				Command:          "./daprd",
				Age:              age.GetAge(d.CreationTimestamp.Time),
				Created:          d.GetCreationTimestamp().String(),
				PID:              -1,
				Replicas:         int(*d.Spec.Replicas),
				SupportsDeletion: false,
				SupportsLogs:     true,
				Address:          fmt.Sprintf("%s-dapr:80", id),
				Status:           fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, d.Status.Replicas),
				Labels:           "app:" + d.Labels["app"],
				Selector:         "app:" + d.Labels["app"],
				Config:           d.Spec.Template.Annotations["dapr.io/config"],
			}

			if val, ok := d.Spec.Template.Annotations[daprPortAnnotation]; ok {
				appPort, err := strconv.Atoi(val)
				if err == nil {
					i.AppPort = appPort
				}
			}

			s := json_serializer.NewYAMLSerializer(json_serializer.DefaultMetaFactory, nil, nil)
			buf := new(bytes.Buffer)
			err := s.Encode(&d, buf)
			if err != nil {
				log.Println(err)
				return list
			}

			i.Manifest = buf.String()

			list = append(list, i)
		}
	}
	return list
}

// getStandaloneInstances returns the Dapr instances running in the standalone environment
func (i *instances) getStandaloneInstances(scope string) []Instance {
	list := []Instance{}
	output, err := standalone.List()
	if err != nil {
		log.Println(err)
	} else {
		for _, o := range output {
			if o.AppID == "" {
				continue
			}
			list = append(list, Instance{
				AppID:            o.AppID,
				HTTPPort:         o.HTTPPort,
				GRPCPort:         o.GRPCPort,
				AppPort:          o.AppPort,
				Command:          o.Command,
				Age:              o.Age,
				Created:          o.Created,
				PID:              o.PID,
				Replicas:         1,
				SupportsDeletion: true,
				SupportsLogs:     false,
				Address:          fmt.Sprintf("localhost:%v", o.HTTPPort),
			})
		}
	}
	return list
}

func (i *instances) getKubernetesScopes() []string {
	scopes := []string{"All"}
	namespaces, err := i.kubeClient.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	if err != nil {
		log.Println(err)
		return scopes
	}
	for _, namespace := range namespaces.Items {
		scopes = append(scopes, namespace.Name)
	}
	return scopes
}

func (i *instances) getStandaloneScopes() []string {
	return []string{"All"}
}
