#!/bin/bash
mkdir -p release/web/
cd web
npm i
ng build
cd ..
cp -r ./web/dist ./release/web/
export GOOS=linux
go build -o dapr-ui
mv ./dapr-ui ./release