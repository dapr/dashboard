package instances

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/dapr/cli/pkg/age"
	"github.com/dapr/cli/pkg/standalone"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

// Instances is an interface to interact with running Dapr instances in Kubrernetes or Standalone modes
type Instances interface {
	Get() []Instance
	Delete(id string) error
	Logs(id string) string
	Configuration(id string) string
}

type instances struct {
	kubeClient     *kubernetes.Clientset
	getInstancesFn func() []Instance
}

// NewInstances returns an Instances implementation
func NewInstances(kubeClient *kubernetes.Clientset) Instances {
	i := instances{}

	if kubeClient != nil {
		i.getInstancesFn = i.getKubernetesInstances
		i.kubeClient = kubeClient
	} else {
		i.getInstancesFn = i.getStandaloneInstances

	}
	return &i
}

func (i *instances) Get() []Instance {
	return i.getInstancesFn()
}

func (i *instances) Logs(id string) string {
	resp, err := i.kubeClient.AppsV1().Deployments(meta_v1.NamespaceAll).List((meta_v1.ListOptions{}))
	if err != nil || len(resp.Items) == 0 {
		return ""
	}

	const daprEnabledAnnotation string = "dapr.io/enabled"
	const daprIDAnnotation string = "dapr.io/id"

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

				for _, p := range pods.Items {
					name := p.ObjectMeta.Name

					options := v1.PodLogOptions{}
					options.Container = "daprd"

					logs := i.kubeClient.CoreV1().Pods(p.ObjectMeta.Namespace).GetLogs(name, &options)
					arr, err := logs.Stream()
					if err != nil {
						log.Println(err)
						return ""
					}

					buf := new(bytes.Buffer)
					buf.ReadFrom(arr)
					logsStr := buf.String()

					return logsStr
				}
			}
		}
	}
	return ""
}

func (i *instances) Configuration(id string) string {
	resp, err := i.kubeClient.AppsV1().Deployments(meta_v1.NamespaceAll).List((meta_v1.ListOptions{}))
	if err != nil || len(resp.Items) == 0 {
		return ""
	}

	const daprEnabledAnnotation string = "dapr.io/enabled"
	const daprIDAnnotation string = "dapr.io/id"

	for _, d := range resp.Items {
		if d.Spec.Template.Annotations[daprEnabledAnnotation] != "" {
			daprID := d.Spec.Template.Annotations[daprIDAnnotation]
			if daprID == id {
				nspace := d.ObjectMeta.Namespace
				restClient := i.kubeClient.CoreV1().RESTClient()
				if err != nil {
					log.Println(err)
					return ""
				}

				url := fmt.Sprintf("/apis/apps/v1/namespaces/%s/deployments/%s", nspace, id)
				data, err := restClient.Get().RequestURI(url).Stream()
				if err != nil {
					log.Println(err)
					return ""
				}

				buf := new(bytes.Buffer)
				buf.ReadFrom(data)
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
	return ""
}

func (i *instances) Delete(id string) error {
	return standalone.Stop(id)
}

func (i *instances) getKubernetesInstances() []Instance {
	list := []Instance{}
	resp, err := i.kubeClient.AppsV1().Deployments(meta_v1.NamespaceAll).List((meta_v1.ListOptions{}))
	if err != nil {
		log.Println(err)
		return list
	}

	const daprEnabledAnnotation string = "dapr.io/enabled"
	const daprIDAnnotation string = "dapr.io/id"
	const daprPortAnnotation string = "dapr.io/id"

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
			}

			if val, ok := d.Spec.Template.Annotations[daprPortAnnotation]; ok {
				appPort, err := strconv.Atoi(val)
				if err == nil {
					i.AppPort = appPort
				}
			}

			s := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)
			buf := new(bytes.Buffer)
			s.Encode(&d, buf)

			i.Manifest = buf.String()

			list = append(list, i)
		}
	}
	return list
}

func (i *instances) getStandaloneInstances() []Instance {
	list := []Instance{}
	output, err := standalone.List()
	if err != nil {
		log.Println(err)
	} else {
		for _, o := range output {
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
