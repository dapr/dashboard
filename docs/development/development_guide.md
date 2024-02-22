# Dapr Dashboard Development Guide

## Environment Setup

## Prerequisites

[Go](https://golang.org/dl/)

[Node & NPM](https://nodejs.org/en/download/)

[Angular CLI](https://angular.io/cli)

[Docker](https://www.docker.com/get-started)

[Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

[Dapr CLI](https://github.com/dapr/cli)

## Expedite UI development

### Build binaries
```sh
make
```

### Run backend

On one terminal, run:

For standalone mode:
```sh
make run-backend-standalone
```

For kubernetes mode:
```sh
make run-backend-kubernetes
```

### Run frontend

On another terminal, run:
```sh
make run-frontend
```

### Run proxy

On a third terminal, run:
```sh
make run-nginx
```

### Open website

Now, open http://localhost:8000. Then, you can change the Angular code and see live changes without rebuilding the website to test every delta.

## Kubernetes

### Build using script
Run: 
```bash
./build_kubernetes.sh <your-username>/<your-image-name>:<your-tag-name> <your-namespace>
```

e.g. `./build_kubernetes.sh docker-username/dashboard:test dapr-system`

### Build from source

#### Build distribution folder and binary
```bash
make
```

#### Build docker image
```bash
docker build -t <your-image-name> .
docker push <your-image-name>
```

#### Apply new deployment
Create a new deployment file with your image name (see [deploy/dashboard.yaml](https://github.com/dapr/dashboard/blob/master/deploy/dashboard.yaml) for reference):
```yaml
...
    spec:
      containers:
      - name: dapr-dashboard
        image: <your-image-name>
        ports:
        - containerPort: 8080
        imagePullPolicy: Always
...
```
Deploy the edited manifest file:
```bash
kubectl apply -f ./test_dashboard.yaml -n <your-namespace>
```
Wait until the dashboard pod is in Running state:
```bash
kubectl get pod --selector=app=dapr-dashboard -w
```
Connect to the dashboard:
```bash
kubectl port-forward svc/dapr-dashboard 8080:8080
```
Alternatively, run:
```bash
dapr dashboard -k
```

Done! point your browser to http://localhost:8080.

## Self-hosted

### Build using script
Run: 
```bash
./build_standalone.sh <your-platform>
```

e.g. `./build_standalone.sh windows_amd64`

### Build from source

#### Build distribution folder and binary
```bash
make

# Windows
./dashboard.exe

# Unix
./dashboard
```

Done! point your browser to http://localhost:8080.

Use the `build.sh` script to generate platform-specific binaries and artifacts.

## Further reference

Check out the other [development guides](https://github.com/dapr/dashboard/tree/master/docs/development), or open an issue with your question. Thank you for contributing to Dapr Dashboard!
