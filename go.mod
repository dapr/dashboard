module github.com/dapr/dashboard

go 1.13

require (
	github.com/dapr/cli v1.0.0-rc.1.0.20201201065415-378e047addfd
	github.com/dapr/dapr v1.0.0-rc.1.0.20201130205314-b9e1e198ebdc
	github.com/gorilla/mux v1.7.3
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.17.8
	k8s.io/apimachinery v0.17.8
	k8s.io/client-go v0.17.2
	k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace k8s.io/client => github.com/kubernetes-client/go v0.0.0-20190928040339-c757968c4c36
