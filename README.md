# Dapr Dashboard

The Dapr Dashboard is a web-based UI for Dapr, allowing users to see information, view logs and more for the Dapr sidecars running either locally or in a Kubernetes cluster.

<p style="text-align:center">
  <img src="img/img.PNG">
</p>

## Features

This repo is under heavy development, all devs and web developers in particular are welcome to start contributing!

### Kubernetes
* List Dapr sidecars with metadata
* View sidecar logs
* View Dapr sidecar deployment file
* View Dapr components in the cluster
* View Dapr configurations in the cluster
* View Dapr control plane status
* View Dapr actor metadata information

### Standalone
* List Dapr sidecars with metadata
* Stop a running sidecar
* View Dapr actor metadata information

## Getting started

### Prerequisites
[Dapr Runtime](https://github.com/dapr/dapr)
[Dapr CLI](https://github.com/dapr/cli)

### Installation

#### Kubernetes
If Dapr was installed with [Helm](https://github.com/dapr/docs/blob/master/getting-started/environment-setup.md#using-helm-advanced), run `dapr dashboard -k`, or if you installed Dapr in a non-default namespace, `dapr dashboard -k -n your-namespace`

If Dapr was installed with `dapr init -k`, run `dapr dashboard -k`

#### Standalone
Running the dashboard locally will work with Dapr instances running on the local machine.

First, make sure you have [Go](https://golang.org/dl/) installed.
Go > 1.13 is required.

Install the Angular CLI:

```
npm i --global @angular/cli
```

Clone the repo and run the dashboard:

```bash
mkdir -p $GOPATH/src/github.com/dapr/dashboard
cd $GOPATH/src/github.com/dapr
git clone git@github.com:dapr/dashboard.git
cd dashboard/web
npm i
ng build
cd ..
go build

# Mac/ Linux
./dashboard

# Windows
./dashboard.exe
```

Done! point your browser to http://localhost:8080.