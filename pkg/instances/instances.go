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
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/dapr/cli/pkg/standalone"
	"github.com/dapr/components-contrib/nameresolution"
	"github.com/dapr/components-contrib/nameresolution/mdns"
	"github.com/dapr/dashboard/pkg/age"
	"github.com/dapr/dashboard/pkg/platforms"
	"github.com/dapr/kit/logger"
	process "github.com/shirou/gopsutil/process"
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
	CheckPlatform() platforms.Platform
}

type instances struct {
	platform          platforms.Platform
	kubeClient        kubernetes.Interface
	getInstancesFn    func(string) []Instance
	getScopesFn       func() []string
	dockerComposePath string
	resolver          nameresolution.Resolver
	metadataClient    http.Client
	daprApiToken      string
}

// NewInstances returns an Instances instance
func NewInstances(platform platforms.Platform, kubeClient *kubernetes.Clientset, dockerComposePath string) Instances {
	i := instances{}
	i.platform = platform
	i.metadataClient = http.Client{}
	i.daprApiToken = os.Getenv("DAPR_API_TOKEN")

	if i.platform == platforms.Kubernetes {
		i.getInstancesFn = i.getKubernetesInstances
		i.getScopesFn = i.getKubernetesScopes
		i.kubeClient = kubeClient
	} else if i.platform == platforms.Standalone {
		i.getInstancesFn = i.getStandaloneInstances
		i.getScopesFn = i.getStandaloneScopes
	} else if i.platform == platforms.DockerCompose {
		i.getInstancesFn = i.getDockerComposeInstances
		i.getScopesFn = i.getDockerComposeScopes
		i.dockerComposePath = dockerComposePath
		i.resolver = mdns.NewResolver(logger.NewLogger("mdns"))
	}
	return &i
}

// Supported checks if the current platform supports Dapr instances
func (i *instances) Supported() bool {
	return i.platform == platforms.Kubernetes || i.platform == platforms.Standalone || i.platform == platforms.DockerCompose
}

// GetScopes returns the result of the appropriate environment's GetScopes function
func (i *instances) GetScopes() []string {
	return i.getScopesFn()
}

// CheckPlatform returns the current environment dashboard is running in
func (i *instances) CheckPlatform() platforms.Platform {
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

// GetInstance uses the appropriate getInstance function (kubernetes, standalone, docker-compose etc.) and returns the given instance from its id
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

	} else if i.platform == platforms.DockerCompose {
		appId := i.GetInstance(scope, id).AppID
		port := i.GetInstance(scope, id).HTTPPort
		address, err := i.resolver.ResolveID(nameresolution.ResolveRequest{ID: appId})
		if err != nil {
			log.Println(err)
			return MetadataOutput{}
		}
		url = append(url, fmt.Sprintf("http://%s:%v/v1.0/metadata", strings.Split(address, ":")[0], port))

	} else {
		port := i.GetInstance(scope, id).HTTPPort
		url = append(url, fmt.Sprintf("http://localhost:%v/v1.0/metadata", port))
	}
	if len(url) != 0 {
		data := i.getMetadataOutputFromURLs(url[0], "")

		if len(url) > 1 {
			// merge the actor metadata from the other replicas

			for index := range url[1:] {
				replicaData := i.getMetadataOutputFromURLs(url[index+1], secondaryUrl[index+1])

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

func (i *instances) getMetadataOutputFromURLs(primaryURL string, secondaryURL string) MetadataOutput {
	req, err := http.NewRequest("GET", primaryURL, nil)

	if len(i.daprApiToken) > 0 {
		req.Header.Add("dapr-api-token", i.daprApiToken)
	}

	resp, err := i.metadataClient.Do(req)
	if err != nil && len(secondaryURL) != 0 {
		log.Println(err)
		secondaryReq, err := http.NewRequest("GET", secondaryURL, nil)
		if len(i.daprApiToken) > 0 {
			secondaryReq.Header.Add("dapr-api-token", i.daprApiToken)
		}
		resp, err = i.metadataClient.Do(secondaryReq)

		if err != nil {
			log.Println(err)
			return MetadataOutput{}
		}
	}

	if err != nil {
		log.Println(err)
		return MetadataOutput{}
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

// getDockerComposeInstances returns the Dapr instances running in the docker compose environment
func (i *instances) getDockerComposeInstances(scope string) []Instance {
	list := []Instance{}

	composeFile, err := os.ReadFile(i.dockerComposePath)
	if err != nil {
		log.Println(err)
		return list
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return list
	}

	configFiles := []types.ConfigFile{}
	configFiles = append(configFiles, types.ConfigFile{
		Filename: "docker-compose.yml",
		Content:  composeFile,
	})

	configDetails := types.ConfigDetails{
		WorkingDir:  workingDir,
		ConfigFiles: configFiles,
		Environment: nil,
	}

	// Can't get creation time and age of daprd services so approximate from dashboard process
	dashboardPID := os.Getpid()
	procDetails, err := process.NewProcess(int32(dashboardPID))
	if err != nil {
		log.Println(err)
		return list
	}

	createUnixTimeMilliseconds, err := procDetails.CreateTime()
	if err != nil {
		log.Println(err)
		return list
	}

	createTime := time.Unix(createUnixTimeMilliseconds/1000, 0)

	project, err := loader.Load(configDetails)

	if err != nil {
		log.Println(err)
	} else {
		for _, service := range project.Services {
			if !strings.Contains(service.Image, "daprd") {
				continue
			}

			appId := ""
			appPort := 80
			foundAppId := false
			foundAppPort := false
			for _, cmdArg := range service.Command {
				if strings.Contains(cmdArg, "app-id") {
					foundAppId = true
					continue
				}

				if strings.Contains(cmdArg, "app-port") {
					foundAppPort = true
					continue
				}

				if foundAppId {
					appId = cmdArg
					foundAppId = false
					continue
				}

				if foundAppPort {
					parsedPort, err := strconv.ParseInt(cmdArg, 10, 0)
					if err == nil {
						appPort = int(parsedPort)
					}
					foundAppId = false
					continue
				}
			}

			list = append(list, Instance{
				AppID:            appId,
				HTTPPort:         3500,
				GRPCPort:         50001,
				AppPort:          appPort,
				Command:          "",
				Age:              age.GetAge(createTime),
				Created:          createTime.Format("2006-01-02 15:04.05"),
				PID:              -1,
				Replicas:         1,
				SupportsDeletion: false,
				SupportsLogs:     false,
				Address:          fmt.Sprintf("%s:3500", appId),
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

func (i *instances) getDockerComposeScopes() []string {
	return []string{"All"}
}

func getAppLabelValue(value string) string {
	if value != "" {
		return "app:" + value
	}

	return ""
}
