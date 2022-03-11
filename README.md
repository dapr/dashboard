# Dapr Dashboard

[![codecov](https://codecov.io/gh/dapr/dashboard/branch/master/graph/badge.svg)](https://codecov.io/gh/dapr/dashboard)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B162%2Fgithub.com%2Fdapr%2Fdashboard.svg?type=shield)](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fdapr%2Fdashboard?ref=badge_shield)

Dapr Dashboard is a web-based UI for Dapr, allowing users to see information, view logs and more for the Dapr applications, components, and configurations running either locally or in a Kubernetes cluster.

<p style="text-align:center">
  <img src="img/img.PNG">
</p>

## Features

Dapr Dashboard provides information about Dapr applications, components, configurations, and control plane services. Users can view metadata, manifests and deployment files, actors, logs, and more on both Kubernetes and self-hosted platforms. For more information, check out the [changelog](docs/development/changelog.md).

## Getting started

### Prerequisites
[Dapr Runtime](https://github.com/dapr/dapr)

[Dapr CLI](https://github.com/dapr/cli)

### Installation

Dapr Dashboard comes pre-packaged with the Dapr CLI. To learn more about the dashboard command, use the CLI command `dapr dashboard -h`.

#### Kubernetes
Run `dapr dashboard -k`, or if you installed Dapr in a non-default namespace, `dapr dashboard -k -n <your-namespace>`.

#### Standalone
Run `dapr dashboard`, and navigate to http://localhost:8080.

### Contributing
Anyone is free to open an issue, a feature request, or a pull request.

To get started in contributing, check out the [development documentation](docs/development/development_guide.md).
