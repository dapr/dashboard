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
	"fmt"
	"testing"
	"time"

	"github.com/dapr/dashboard/pkg/platforms"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func newTestSimpleK8s(objects ...runtime.Object) instances {
	client := instances{}
	client.kubeClient = fake.NewSimpleClientset(objects...) //nolint:staticcheck
	return client
}

func newDaprControlPlanePod(name string, appName string, creationTime time.Time, state v1.ContainerState, ready bool) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "dapr-system",
			Annotations: map[string]string{},
			Labels: map[string]string{
				"app": appName,
			},
			CreationTimestamp: metav1.Time{
				Time: creationTime,
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					State: state,
					Ready: ready,
				},
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image: name + ":0.0.1",
				},
			},
		},
	}
}

func TestControlPlaneServices(t *testing.T) {
	controlPlaneServices := []struct {
		name    string
		appName string
	}{
		{"dapr-operator-67d7d7bb6c-7h96c", "dapr-operator"},
		{"dapr-operator-67d7d7bb6c-2h96d", "dapr-operator"},
		{"dapr-operator-67d7d7bb6c-3h96c", "dapr-operator"},
		{"dapr-placement-server-0", "dapr-placement-server"},
		{"dapr-placement-server-1", "dapr-placement-server"},
		{"dapr-placement-server-2", "dapr-placement-server"},
		{"dapr-sentry-647759cd46-9ptks", "dapr-sentry"},
		{"dapr-sentry-647759cd46-aptks", "dapr-sentry"},
		{"dapr-sentry-647759cd46-bptks", "dapr-sentry"},
		{"dapr-sidecar-injector-74648c9dcb-5bsmn", "dapr-sidecar-injector"},
		{"dapr-sidecar-injector-74648c9dcb-6bsmn", "dapr-sidecar-injector"},
		{"dapr-sidecar-injector-74648c9dcb-7bsmn", "dapr-sidecar-injector"},
	}

	runtimeObj := make([]runtime.Object, len(controlPlaneServices))
	for i, s := range controlPlaneServices {
		testTime := time.Now()
		runtimeObj[i] = newDaprControlPlanePod(
			s.name, s.appName,
			testTime.Add(time.Duration(-20)*time.Minute),
			v1.ContainerState{
				Running: &v1.ContainerStateRunning{
					StartedAt: metav1.Time{
						Time: testTime.Add(time.Duration(-19) * time.Minute),
					},
				},
			}, true)
	}

	k8s := newTestSimpleK8s(runtimeObj...)
	status := k8s.GetControlPlaneStatus()
	assert.Equal(t, 12, len(status), "Expected status list length to match")
}

func TestSupported(t *testing.T) {
	var scenarios = []struct {
		platform platforms.Platform
		want     bool
	}{
		{platforms.Kubernetes, true},
		{platforms.Standalone, true},
		{platforms.DockerCompose, true},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Platform %s should be supported", scenario.platform), func(t *testing.T) {
			target := NewInstances(scenario.platform, nil, "")
			isSupported := target.Supported()
			assert.Equal(t, scenario.want, isSupported)
		})
	}
}

func TestDockerComposeGetInstances(t *testing.T) {
	target := NewInstances(platforms.DockerCompose, nil, "testdata/docker-compose.yml")
	instances := target.GetInstances("")
	assert.Equal(t, 1, len(instances), "Should parse docker compose file and detect one instance")
}

func TestDockerComposeGetInstance(t *testing.T) {
	var scenarios = []struct {
		id   string
		want bool
	}{
		{"MyApplication.DaprSidecar", true},
		{"does_not_exist", false},
	}

	target := NewInstances(platforms.DockerCompose, nil, "testdata/docker-compose.yml")

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("Should load valid instance data - %t", scenario.want), func(t *testing.T) {
			instance := target.GetInstance("", scenario.id)
			assert.NotNil(t, instance, "Should always return something")

			if scenario.want {
				assert.Equal(t, scenario.id, instance.AppID, "Should return the correct instance")
				assert.Equal(t, 3500, instance.HTTPPort, "Port should be set")
				assert.Equal(t, false, instance.SupportsLogs, "Logs are not supported")
				assert.Equal(t, false, instance.SupportsDeletion, "Delegation is not supported")
				assert.Equal(t, "MyApplication.DaprSidecar:3500", instance.Address, "Address should be set")
				assert.Equal(t, 80, instance.AppPort, "AppPort should be set")
			} else {
				assert.Empty(t, instance.AppID, "When instance not valid, AppID is not set")
			}
		})
	}
}
