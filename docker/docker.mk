# ------------------------------------------------------------
# Copyright 2021 The Dapr Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------


# Docker image build and push setting
DOCKER:=docker
DOCKERFILE_DIR?=./docker

DAPR_SYSTEM_IMAGE_NAME=$(RELEASE_NAME)
DAPR_RUNTIME_IMAGE_NAME=dashboard

# build docker image for linux
BIN_PATH=$(OUT_DIR)/$(TARGET_OS)_$(TARGET_ARCH)

DOCKERFILE:=Dockerfile
ifeq ($(TARGET_OS), windows)
  DOCKERFILE:=Dockerfile-windows
endif

# Supported docker image architecture
DOCKERMUTI_ARCH=linux-amd64 linux-arm linux-arm64 windows-amd64

################################################################################
# Target: docker-build, docker-push                                            #
################################################################################

LINUX_BINS_OUT_DIR=$(OUT_DIR)/linux_$(GOARCH)
DOCKER_IMAGE_TAG=$(DAPR_REGISTRY)/$(DAPR_SYSTEM_IMAGE_NAME):$(DAPR_TAG)
DAPR_RUNTIME_DOCKER_IMAGE_TAG=$(DAPR_REGISTRY)/$(DAPR_RUNTIME_IMAGE_NAME):$(DAPR_TAG)
DAPR_PLACEMENT_DOCKER_IMAGE_TAG=$(DAPR_REGISTRY)/$(DAPR_PLACEMENT_IMAGE_NAME):$(DAPR_TAG)
DAPR_SENTRY_DOCKER_IMAGE_TAG=$(DAPR_REGISTRY)/$(DAPR_SENTRY_IMAGE_NAME):$(DAPR_TAG)

ifeq ($(LATEST_RELEASE),true)
DOCKER_IMAGE_LATEST_TAG=$(DAPR_REGISTRY)/$(DAPR_SYSTEM_IMAGE_NAME):$(LATEST_TAG)
DAPR_RUNTIME_DOCKER_IMAGE_LATEST_TAG=$(DAPR_REGISTRY)/$(DAPR_RUNTIME_IMAGE_NAME):$(LATEST_TAG)
DAPR_PLACEMENT_DOCKER_IMAGE_LATEST_TAG=$(DAPR_REGISTRY)/$(DAPR_PLACEMENT_IMAGE_NAME):$(LATEST_TAG)
DAPR_SENTRY_DOCKER_IMAGE_LATEST_TAG=$(DAPR_REGISTRY)/$(DAPR_SENTRY_IMAGE_NAME):$(LATEST_TAG)
endif


# To use buildx: https://github.com/docker/buildx#docker-ce
export DOCKER_CLI_EXPERIMENTAL=enabled

# check the required environment variables
check-docker-env:
ifeq ($(DAPR_REGISTRY),)
	$(error DAPR_REGISTRY environment variable must be set)
endif
ifeq ($(DAPR_TAG),)
	$(error DAPR_TAG environment variable must be set)
endif

check-arch:
ifeq ($(TARGET_OS),)
	$(error TARGET_OS environment variable must be set)
endif
ifeq ($(TARGET_ARCH),)
	$(error TARGET_ARCH environment variable must be set)
endif


docker-build: check-docker-env check-arch
	$(info Building $(DAPR_RUNTIME_DOCKER_IMAGE_TAG) docker image ...)
	$(DOCKER) build --build-arg BIN_PATH=$(BIN_PATH) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) . -t $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH) --platform $(TARGET_OS)/$(TARGET_ARCH)

# push docker image to the registry
docker-push: docker-build
	$(info Pushing $(DAPR_RUNTIME_DOCKER_IMAGE_TAG) docker image ...)
	$(DOCKER) push $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)

# publish muti-arch docker image to the registry
docker-manifest-create: check-docker-env
	$(DOCKER) manifest create $(DAPR_RUNTIME_DOCKER_IMAGE_TAG) $(DOCKERMUTI_ARCH:%=$(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-%)
ifeq ($(LATEST_RELEASE),true)
	$(DOCKER) manifest create $(DAPR_RUNTIME_DOCKER_IMAGE_LATEST_TAG) $(DOCKERMUTI_ARCH:%=$(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-%)
endif

docker-publish: docker-manifest-create
	$(DOCKER) manifest push $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)
ifeq ($(LATEST_RELEASE),true)
	$(DOCKER) manifest push $(DAPR_RUNTIME_DOCKER_IMAGE_LATEST_TAG)
endif

check-windows-version:
ifeq ($(WINDOWS_VERSION),)
	$(error WINDOWS_VERSION environment variable must be set)
endif

docker-windows-base-build: check-windows-version
	$(DOCKER) build --build-arg WINDOWS_VERSION=$(WINDOWS_VERSION) -f $(DOCKERFILE_DIR)/$(DOCKERFILE)-base . -t $(DAPR_REGISTRY)/windows-base:$(WINDOWS_VERSION)

docker-windows-base-push: check-windows-version
	$(DOCKER) push $(DAPR_REGISTRY)/windows-base:$(WINDOWS_VERSION)