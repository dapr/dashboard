#!/bin/bash
mkdir -p release/web/
cd web
npm i
ng build --base-href ./
cd ..
cp -r ./web/dist ./release/web/
export GOOS=linux
go build -o dashboard.exe
mv ./dashboard ./release