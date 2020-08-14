#!/bin/bash

# USAGE: ./build_kubernetes.sh <your_image> <your_namespace>
# e.g. ./build_kubernetes.sh docker-username/dashboard:test dapr-system

cd web
npm i
ng build
cd ..
release_dir=./release
artifacts_dir=${release_dir}/artifacts
rm -r -f ${release_dir}
mkdir -p ${artifacts_dir}

GOOS="linux"
GOARCH="amd64"
platform="linux_amd64"
image="$1"
namespace="$2"
  
go_executable_output_file=dist/${platform}/release/dashboard

echo building go executable for $GOOS $GOARCH, output will be $go_executable_output_file
env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a -o $go_executable_output_file

platform_release_dir=${release_dir}/${platform}

echo preparing release dir ${platform_release_dir}
mkdir -p ${platform_release_dir}/web/
cp -r ./web/dist ${platform_release_dir}/web/
mv ./$go_executable_output_file ${platform_release_dir}

docker build -t ${image} .
docker push ${image}

kubectl delete -f ./test_dashboard.yaml -n ${namespace}

kubectl apply -f ./test_dashboard.yaml -n ${namespace}

sleep 5

dapr dashboard -k
