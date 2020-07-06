#!/bin/bash
# build node angular code
cd web
npm i
ng build --base-href ./
cd ..
release_dir=./release
artifacts_dir=${release_dir}/artifacts
rm -r -f ${release_dir}
mkdir -p ${artifacts_dir}


# build go code and prepare release dir for platforms.
# (TODO: currently its for linux enable darwin and windows (if angular code build here can be used on windows and mac) by uncommenting the platforms array with multiple platforms below)

platforms=("linux_amd64")
#platforms=("linux_amd64" "windows_amd64" "darwin_amd64")

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

  platform_release_dir=${release_dir}/${platform}
  platform_artifact_archive=${artifacts_dir}/dashboard_${platform}

  echo preparing release dir ${platform_release_dir}
  mkdir -p ${platform_release_dir}/web/
  cp -r ./web/dist ${platform_release_dir}/web/
  mv ./$go_executable_output_file ${platform_release_dir}

  # create archives
  tar -zcf ${platform_artifact_archive}.tar.gz ${platform_release_dir}
done
