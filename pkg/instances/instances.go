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

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

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
	daprIDAnnotation      = "dapr.io/app-id"
	daprPortAnnotation    = "dapr.io/port"
)

var controlPlaneLabels = [...]string{
	"dapr-operator",
	"dapr-sentry",
	"dapr-placement",
	"dapr-placement-server",
	"dapr-sidecar-injector",
	"dapr-dashboard",
}

// Instances is an interface to interact with running Dapr instances in Kubernetes or Standalone modes
type Instances interface {
	Supported() bool
	GetInstances(scope string) []Instance
	GetInstance(scope string, id string) Instance
	DeleteInstance(scope string, id string) error
	GetContainers(scope string, id string) []string
	GetLogStream(scope, id, containerName string) ([]io.ReadCloser, error)
	GetDeploymentConfiguration(scope string, id string) string
	GetControlPlaneStatus() []StatusOutput
	GetMetadata(scope string, id string) MetadataOutput
	GetScopes() []string
	CheckPlatform() string
}

type instances struct {
	platform       string
	kubeClient     kubernetes.Interface
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

// GetContainers returns a list of containers for an app.
func (i *instances) GetContainers(scope string, id string) []string {
	ctx := context.Background()
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List(ctx, (meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return nil
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(ctx, meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						log.Println(err)
						return nil
					}

					if len(pods.Items) > 0 {
						p := pods.Items[0]
						out := []string{}
						for _, container := range p.Spec.Containers {
							out = append(out, container.Name)
						}
						return out
					}
				}
			}
		}
	}
	return nil
}

// GetLogStream returns a stream of bytes from K8s logs
func (i *instances) GetLogStream(scope, id, containerName string) ([]io.ReadCloser, error) {
	ctx := context.Background()
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List(ctx, (meta_v1.ListOptions{}))
		if err != nil {
			return nil, err
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(ctx, meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						return nil, err
					}

					var logstreams []io.ReadCloser

					for _, p := range pods.Items {
						name := p.ObjectMeta.Name

						for _, container := range p.Spec.Containers {
							if container.Name == containerName {
								var tailLines int64 = 100
								options := v1.PodLogOptions{}
								options.Container = container.Name
								options.Timestamps = true
								options.TailLines = &tailLines
								if len(pods.Items) == 1 {
									options.Follow = true
								} else {
									options.Follow = false // this is necessary to show logs from multiple replicas
								}

								res := i.kubeClient.CoreV1().Pods(p.ObjectMeta.Namespace).GetLogs(name, &options)
								stream, streamErr := res.Stream(ctx)
								if streamErr == nil {
									logstreams = append(logstreams, stream)
								}
							}
						}
					}
					return logstreams, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("could not find logstream for %v, %v, %v", scope, id, containerName)
}

// GetDeploymentConfiguration returns the metadata of a Dapr application in YAML format
func (i *instances) GetDeploymentConfiguration(scope string, id string) string {
	ctx := context.Background()
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List(ctx, (meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return ""
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(ctx, meta_v1.ListOptions{
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
						data, err := restClient.Get().RequestURI(url).Stream(ctx)
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
	apps, err := standalone.List()
	if err != nil {
		return err
	}
	cliPIDToNoOfApps := standalone.GetCLIPIDCountMap(apps)
	return standalone.Stop(id, cliPIDToNoOfApps, apps)
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
	ctx := context.Background()
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

				p, err := i.kubeClient.CoreV1().Pods(v1.NamespaceAll).List(ctx, options)
				if err == nil {
					for _, pod := range p.Items {
						name := pod.Name
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
							Service:   label,
							Name:      name,
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
				}
				wg.Done()
			}(lbl)
		}

		wg.Wait()
		return statuses
	}
	return []StatusOutput{}
}

// GetMetadata returns the result from the /v1.0/metadata endpoint from all replicas
func (i *instances) GetMetadata(scope string, id string) MetadataOutput {
	ctx := context.Background()
	var url []string
	var secondaryUrl []string
	if i.kubeClient != nil {
		resp, err := i.kubeClient.AppsV1().Deployments(scope).List(ctx, (meta_v1.ListOptions{}))
		if err != nil || len(resp.Items) == 0 {
			return MetadataOutput{}
		}

		for _, d := range resp.Items {
			if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
				daprID := d.Spec.Template.Annotations[daprIDAnnotation]
				if daprID == id {
					pods, err := i.kubeClient.CoreV1().Pods(d.GetNamespace()).List(ctx, meta_v1.ListOptions{
						LabelSelector: labels.SelectorFromSet(d.Spec.Selector.MatchLabels).String(),
					})
					if err != nil {
						log.Println(err)
						return MetadataOutput{}
					}

					if len(pods.Items) > 0 {
						p := pods.Items[0]
						url = append(url, fmt.Sprintf("http://%v:%v/v1.0/metadata", p.Status.PodIP, 3501))
						secondaryUrl = append(secondaryUrl, fmt.Sprintf("http://%v:%v/v1.0/metadata", p.Status.PodIP, 3500))
					}
				}
			}
		}

	} else {
		port := i.GetInstance(scope, id).HTTPPort
		url = append(url, fmt.Sprintf("http://localhost:%v/v1.0/metadata", port))
	}
	if len(url) != 0 {
		data := getMetadataOutputFromURLs(url[0], "")
		if len(url) > 1 {
			// merge the actor metadata from the other replicas

			for i := range url[1:] {
				replicaData := getMetadataOutputFromURLs(url[i+1], secondaryUrl[i+1])

				for _, actor := range replicaData.Actors {
					// check if this actor type is already in the list
					found := false

					for _, knownActor := range data.Actors {
						if knownActor.Type == actor.Type {
							found = true
							knownActor.Count += actor.Count
							break
						}
					}

					if !found {
						data.Actors = append(data.Actors, actor)
					}
				}
			}
		}

		return data
	}
	return MetadataOutput{}
}

func getMetadataOutputFromURLs(primaryURL string, secondaryURL string) MetadataOutput {
	resp, err := http.Get(primaryURL)
	if err != nil && len(secondaryURL) != 0 {
		log.Println(err)
		resp, err = http.Get(secondaryURL)
		if err != nil {
			log.Println(err)
			return MetadataOutput{}
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

// GetInstances returns the result of the appropriate environment's GetInstance function
func (i *instances) GetInstances(scope string) []Instance {
	return i.getInstancesFn(scope)
}

// getKubernetesInstances gets the list of Dapr applications running in the Kubernetes environment
func (i *instances) getKubernetesInstances(scope string) []Instance {
	ctx := context.Background()
	list := []Instance{}
	resp, err := i.kubeClient.AppsV1().Deployments(scope).List(ctx, (meta_v1.ListOptions{}))
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
				Command:          "",
				Age:              age.GetAge(d.CreationTimestamp.Time),
				Created:          d.GetCreationTimestamp().String(),
				PID:              -1,
				Replicas:         int(*d.Spec.Replicas),
				SupportsDeletion: false,
				SupportsLogs:     true,
				Address:          fmt.Sprintf("%s-dapr:80", id),
				Status:           fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, d.Status.Replicas),
				Labels:           getAppLabelValue(d.Labels["app"]),
				Selector:         getAppLabelValue(d.Spec.Selector.MatchLabels["app"]),
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
				PID:              o.DaprdPID,
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
	ctx := context.Background()
	scopes := []string{"All"}
	namespaces, err := i.kubeClient.CoreV1().Namespaces().List(ctx, meta_v1.ListOptions{})
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

func getAppLabelValue(value string) string {
	if value != "" {
		return "app:" + value
	}

	return ""
}
