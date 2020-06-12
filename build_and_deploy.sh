#!/bin/bash
mkdir -p release/web/
cd web
npm i
ng build
cd ..
cp -r ./web/dist ./release/web/
export GOOS=linux
go build -o dashboard
mv ./dashboard ./release

docker build -t willdavsmith/dashboard .
docker push willdavsmith/dashboard

kubectl get pods --no-headers=true | awk '/dapr-dashboard/{print $1}' | xargs kubectl delete pod
kubectl apply -f ./deploy/dashboard.yaml

sleep 5

kubectl port-forward svc/dapr-dashboard 8080:8080