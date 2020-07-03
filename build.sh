#!/bin/bash
# build node angular code
cd web
npm i
ng build --base-href ./
cd ..
rm -r -f release

# build go code and prepare release dir for platforms.
# (TODO: currently its for linux enable darwin and windows as well by uncommenting the platforms array with multiple platforms below)

platforms=("linux_amd64")
platforms=("linux_amd64" "windows_amd64" "darwin_amd64")

for platform in "${platforms[@]}"
do
  platform_split=(${platform//_/ })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  
  go_executable_output_file=dist/${platform}/release/dashboard
  if [ $GOOS = "windows" ]; then
    go_executable_output_file+='.exe'
  fi

  echo building go executable for $GOOS $GOARCH, output will be $go_executable_output_file
  env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a -o $go_executable_output_file

  echo preparing release dir release/${platform}
  mkdir -p ./release/${platform}/web/
  cp -r ./web/dist ./release/${platform}/web/
  mv ./$go_executable_output_file ./release/${platform}

  # create archives
  if [ $GOOS = "windows" ]; then
    zip -r -q ./release/artifactes/dashboard_${platform}.zip ./release/${platform}
  else
    tar -zcf ./release/artifacts/dashboard_${platform}.tar.gz ./release/${platform}
  fi
done
