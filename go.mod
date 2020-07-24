module github.com/dapr/dashboard

go 1.13

require (
	github.com/dapr/cli v0.9.0
	github.com/dapr/dapr v0.9.0
	github.com/dapr/go-sdk v0.0.0-20200312165010-7bb7a2205f3b // indirect
	github.com/gorilla/mux v1.7.3
	github.com/joomcode/errorx v1.0.1 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace k8s.io/client => github.com/kubernetes-client/go v0.0.0-20190928040339-c757968c4c36
