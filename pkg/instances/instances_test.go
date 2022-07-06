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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func newTestSimpleK8s(objects ...runtime.Object) instances {
	client := instances{}
	client.kubeClient = fake.NewSimpleClientset(objects...)
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
