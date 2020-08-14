#!/bin/bash
# build node angular code

# --generate-artifacts flag is passed in github actions workflow to generate artifacts.
GENERATE_ARTIFACTS=false
RELEASE_DIR=./release
DASHBOARD_VERSION=edge
OTHER_ARGUMENTS=()

for arg in "$@"
do
    case $arg in
        -a|--generate-artifacts)
        GENERATE_ARTIFACTS=true
        shift # Remove ----generate-artifacts from processing
        ;;
        -v|--release-version)
        DASHBOARD_VERSION="$2"
        shift # Remove argument name from processing
        shift # Remove argument value from processing
        ;;
        -r|--release-dir)
        RELEASE_DIR="$2"
        shift # Remove argument name from processing
        shift # Remove argument value from processing
        ;;
        *)
        OTHER_ARGUMENTS+=("$1")
        shift # Remove generic argument from processing
        ;;
    esac
done

echo "# Release directory: $RELEASE_DIR"
echo "# Generate artifacts: $GENERATE_ARTIFACTS"
echo "# Version: $DASHBOARD_VERSION"
#echo "# Other arguments: ${OTHER_ARGUMENTS[*]}"

cd web
npm i
ng build
cd ..
rm -r -f ${RELEASE_DIR}
mkdir -p ${RELEASE_DIR}


# build go code and prepare release dir for platforms.
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
  env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -a -o $go_executable_output_file -ldflags="-X 'github.com/dapr/dashboard/pkg/version.Version=$DASHBOARD_VERSION'"

  platform_release_dir=${RELEASE_DIR}/${platform}  

  echo preparing release dir ${platform_release_dir}
  mkdir -p ${platform_release_dir}/web/
  cp -r ./web/dist ${platform_release_dir}/web/
  mv ./$go_executable_output_file ${platform_release_dir}

  # create archives if specified
  if [ $GENERATE_ARTIFACTS = "true" ]; then      
    artifacts_dir=${RELEASE_DIR}/artifacts
    mkdir -p ${artifacts_dir}
    platform_artifact_archive=${artifacts_dir}/dashboard_${platform}
    echo generating artifacts in ${artifacts_dir}
    
    tar -zcf ${platform_artifact_archive}.tar.gz ${platform_release_dir}
    
    # generate .zip also (if zip is installed)
    if hash zip 2>/dev/null; then
      zip -r -q ${platform_artifact_archive}.zip ${platform_release_dir}
    else
      echo skipping generation of ${platform_artifact_archive}.zip as zip is not installed. 
    fi
  fi
done