# Development Guide

## Prerequisites
[Go > 1.13](https://golang.org/dl/)

[Angular CLI](https://angular.io/cli)

[Docker](https://www.docker.com/get-started)

[Kubernetes](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Building from source

Build the website and the Go web server into the `/release` dir:

```
./build.sh
```

## Building a Docker image

After you have the release dir in place, run:

```
docker build -t <image>:<tag> .
docker push <image>:<tag>
```

Edit deploy/dashboard.yaml to reference your docker image:
```...
    spec:
      containers:
      - name: dapr-dashboard
        image: <image>:<tag>
        ports:
        - containerPort: 8080
        imagePullPolicy: Always
```
Deploy the edited manifest file:
```
kubectl apply -f ./deploy/dashboard.yaml
```
Wait until the dashboard pod is in Running state:
```
kubectl get pod --selector=app=dapr-dashboard -w
```
Connect to the dashboard:
```
kubectl port-forward svc/dapr-dashboard 8080:8080
```
Done! point your browser to http://localhost:8080.