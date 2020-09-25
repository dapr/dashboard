# ------------------------------------------------------------
# Copyright (c) Microsoft Corporation.
# Licensed under the MIT License.
# ------------------------------------------------------------

################################################################################
# Variables                                                                    #
################################################################################

export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0
BINARIES    ?= dashboard

# Add latest tag if LATEST_RELEASE is true
LATEST_RELEASE ?=

ifdef REL_VERSION
	DAPR_VERSION := $(REL_VERSION)
else
	DAPR_VERSION := edge
endif

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

ifeq ($(GOARCH),amd64)
	LATEST_TAG=latest
else
	LATEST_TAG=latest-$(GOARCH)
endif

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
else
   TARGET_OS_LOCAL ?= windows
endif
export GOOS ?= $(TARGET_OS_LOCAL)

# Default docker container and e2e test targst.
TARGET_OS ?= linux
TARGET_ARCH ?= amd64

ifeq ($(GOOS),windows)
BINARY_EXT_LOCAL:=.exe
GOLANGCI_LINT:=golangci-lint.exe
else
BINARY_EXT_LOCAL:=
GOLANGCI_LINT:=golangci-lint
endif

export BINARY_EXT ?= $(BINARY_EXT_LOCAL)

OUT_DIR := ./release

################################################################################
# Go build details                                                             #
################################################################################
BASE_PACKAGE_NAME := github.com/dapr/dashboard

DEFAULT_LDFLAGS:=-X $(BASE_PACKAGE_NAME)/pkg/version.commit=$(GIT_VERSION) -X $(BASE_PACKAGE_NAME)/pkg/version.version=$(DAPR_VERSION)

ifeq ($(origin DEBUG), undefined)
  OUT_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else ifeq ($(DEBUG),0)
  OUT_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else
  OUT_DIR:=debug
  GCFLAGS:=-gcflags="all=-N -l"
  LDFLAGS:="$(DEFAULT_LDFLAGS)"
  $(info Build with debugger information)
endif

DAPR_ARTIFACTS := $(OUT_DIR)/artifacts
DAPR_OUT_DIR := $(OUT_DIR)/$(GOOS)_$(GOARCH)
DAPR_LINUX_OUT_DIR := $(OUT_DIR)/linux_$(GOARCH)

################################################################################
# Target: build                                                                #
################################################################################
.PHONY: build
DAPR_BINS:=$(foreach ITEM,$(BINARIES),$(DAPR_OUT_DIR)/$(ITEM)$(BINARY_EXT))
build: $(DAPR_BINS)

# Generate builds for dapr binaries for the target
# Params:
# $(1): the binary name for the target
# $(2): the binary main directory
# $(3): the target os
# $(4): the target arch
# $(5): the output directory
define genBinariesForTarget
.PHONY: $(5)/$(1)
$(5)/$(1):
	CGO_ENABLED=$(CGO) GOOS=$(3) GOARCH=$(4) go build $(GCFLAGS) -ldflags=$(LDFLAGS) \
	-o $(5)/$(1) \
	./main.go;
	mkdir -p $(5)/web/dist
	cd web && npm i && ng build --outputPath=../$(5)/web/dist
endef

# Generate binary targets
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM)$(BINARY_EXT),./cmd/$(ITEM),$(GOOS),$(GOARCH),$(DAPR_OUT_DIR))))

################################################################################
# Target: build-linux                                                          #
################################################################################
BUILD_LINUX_BINS:=$(foreach ITEM,$(BINARIES),$(DAPR_LINUX_OUT_DIR)/$(ITEM))
build-linux: $(BUILD_LINUX_BINS)

# Generate linux binaries targets to build linux docker image
ifneq ($(GOOS), linux)
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM),./cmd/$(ITEM),linux,$(GOARCH),$(DAPR_LINUX_OUT_DIR))))
endif

################################################################################
# Target: archive                                                              #
################################################################################
ARCHIVE_OUT_DIR ?= $(DAPR_OUT_DIR)
ARCHIVE_FILE_ZIPS:=$(foreach ITEM,$(BINARIES),dashboard_$(ITEM).zip)
ARCHIVE_FILE_TGZS:=$(foreach ITEM,$(BINARIES),dashboard_$(ITEM).tar.gz)

archive: $(ARCHIVE_FILE_ZIPS) $(ARCHIVE_FILE_TGZS)

# Generate archive files for each binary
# $(1): the binary name to be archived
# $(2): the archived file output directory
define genArchiveBinary
ifeq ($(GOOS),windows)
dashboard_$(1).zip:
	7z.exe a -tzip "$(ARCHIVE_OUT_DIR)/$(1)_$(GOOS)_$(GOARCH).zip" "$(DAPR_OUT_DIR)"
else
dashboard_$(1).zip:
	zip -r -q "$(ARCHIVE_OUT_DIR)/$(1)_$(GOOS)_$(GOARCH).zip" "$(DAPR_OUT_DIR)"
endif
dashboard_$(1).tar.gz:
	tar -zcf "$(ARCHIVE_OUT_DIR)/$(1)_$(GOOS)_$(GOARCH).tar.gz" "$(DAPR_OUT_DIR)"
endef

# Generate archive-*.[zip|tar.gz] targets
$(foreach ITEM,$(BINARIES),$(eval $(call genArchiveBinary,$(ITEM),$(ARCHIVE_OUT_DIR))))


################################################################################
# Target: archive                                                              #
################################################################################
release: build archive

################################################################################
# Target: test                                                                 #
################################################################################
.PHONY: test
test:
	go test ./pkg/... $(COVERAGE_OPTS)

################################################################################
# Target: lint                                                                 #
################################################################################
# Due to https://github.com/golangci/golangci-lint/issues/580, we need to add --fix for windows
.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --timeout=20m

################################################################################
# Target: docker                                                               #
################################################################################
include docker/docker.mk