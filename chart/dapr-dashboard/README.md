## Install

```sh
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm install dapr-dashboard dapr/dapr-dashboard
```

## Chart Options

| Parameter                               | Description                                                                                                                                                            | Default            |
|-----------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| `replicaCount`           | Number of replicas                                                                                                                                                     | `1`                |
| `logLevel`               | service Log level                                                                                                                                                      | `info`             |
| `image.registry`         | docker registry                                                                                                                                                        | `docker.io/daprio` |
| `image.imagePullPolicy`                  | Global Control plane service imagePullPolicy                            | `IfNotPresent`          |
| `image.imagePullSecrets` | docker images pull secrets for docker registry                                                                                                                         | `docker.io/daprio` |
| `image.name`             | docker image name                                                                                                                                                      | `dashboard`        |
| `image.tag`              | docker image tag                                                                                                                                                       | latest release         |
| `serviceType`            | Type of [Kubernetes service](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) to use for the Dapr Dashboard service | `ClusterIP`        |
| `runAsNonRoot`           | Boolean value for `securityContext.runAsNonRoot`. You may have to set this to `false` when running in Minikube                                                         | `true`             |
| `resources`              | Value of `resources` attribute. Can be used to set memory/cpu resources/limits. See the section "Resource configuration" above. Defaults to empty                      | `{}`               |
| `ingress.enabled`        | Boolean value for enabling the creation of the ingress resource                                                                                                        | `false`            |
| `ingress.className`      | ingress className of the ingress controller (e.g.nginx)                                                                                                                | ``                 |
| `ingress.host`           | Fully qualified hostname of the dashboard URL (e.g `dashboard.dapr.local`) | ``                 |
| `ingress.tls.enabled`    | If true, enables TLS on the ingress for the Dashboard                                                                                                                      | `false`            |
| `ingress.tls.secretName` | Name of the Kubernetes secret containing the TLS certificate (key/certificate) for the Dashboard. Ignored if `dapr_dashboard.ingress.tls.enabled` is `false`. | ``                 |
| `registry`                         | Docker image registry                                                   | `docker.io/daprio`      |
| `tag`                              | Docker image version tag                                                | latest release          |
| `logAsJson`                        | Json log format for control plane services                              | `false`                 |
| `ha.enabled`                       | Highly Availability mode enabled for control plane                      | `false`                 |
| `ha.replicaCount`                  | Number of replicas of control plane services in Highly Availability mode  | `3`                   |
| `ha.disruption.minimumAvailable`   | Minimum amount of available instances for control plane. This can either be effective count or %. | ``             |
| `ha.disruption.maximumUnavailable` | Maximum amount of instances that are allowed to be unavailable for control plane. This can either be effective count or %. | `25%`             |              |
| `daprControlPlaneOs`               | Operating System for Dapr control plane                                 | `linux`                 |
| `daprControlPlaneArch`             | CPU Architecture for Dapr control plane                                 | `amd64`                 |
| `nodeSelector`                     | Pods will be scheduled onto a node node whose labels match the nodeSelector        | `{}`         |
| `tolerations`                      | Pods will be allowed to schedule onto a node whose taints match the tolerations    | `{}`         |
| `labels`                           | Custom pod labels                                                                  | `{}`         |
| `k8sLabels`                        | Custom metadata labels                                                             | `{}`         |
| `rbac.namespaced`                  | Removes cluster wide permissions where applicable  | `false` |