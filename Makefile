VERSION ?= main

# Go architecture and targets images to build
GOARCH ?= amd64
MULTIARCH_TARGETS ?= amd64

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL := /usr/bin/env bash

IMAGE_ORG ?= jotak
IMAGE_NAME ?= quay.io/${IMAGE_ORG}/net-infra-mon
IMAGE ?= ${IMAGE_TAG_BASE}:${VERSION}

# Image building tool (docker / podman) - docker is preferred in CI
OCI_BIN_PATH = $(shell which docker 2>/dev/null || which podman)
OCI_BIN ?= $(shell basename ${OCI_BIN_PATH})
OCI_BUILD_OPTS ?=

ifneq ($(CLEAN_BUILD),)
	BUILD_DATE := $(shell date +%Y-%m-%d\ %H:%M)
	BUILD_SHA := $(shell git rev-parse --short HEAD)
	GO_BUILD_OPTS ?= -ldflags "-X 'main.buildVersion=${VERSION}-${BUILD_SHA}' -X 'main.buildDate=${BUILD_DATE}'"
endif

.DEFAULT_GOAL := help

# build a single arch target provided as argument
define build_target
	echo 'building image for arch $(1)'; \
	DOCKER_BUILDKIT=1 $(OCI_BIN) buildx build --load --build-arg GO_BUILD_OPTS="${GO_BUILD_OPTS}" --build-arg TARGETARCH=$(1) ${OCI_BUILD_OPTS} -t ${IMAGE}-$(1) -f Dockerfile .;
endef

# push a single arch target image
define push_target
	echo 'pushing image ${IMAGE}-$(1)'; \
	DOCKER_BUILDKIT=1 $(OCI_BIN) push ${IMAGE}-$(1);
endef

# manifest create a single arch target provided as argument
define manifest_add_target
	echo 'manifest add target $(1)'; \
	DOCKER_BUILDKIT=1 $(OCI_BIN) manifest add ${IMAGE} ${IMAGE}-$(target);
endef

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: prereqs
prereqs: ## Test if prerequisites are met, and installing missing dependencies
	@echo "### Test if prerequisites are met, and installing missing dependencies"
	GOFLAGS="" go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

.PHONY: vendors
vendors: ## Check go vendors
	@echo "### Checking vendors"
	go mod tidy && go mod vendor

##@ Build

.PHONY: build
build: fmt ## Build go binary
	@echo "### Building go binary"
	GOARCH=${GOARCH} go build ${GO_BUILD_OPTS} -mod vendor -o nim-monitoring-agent cmd/agent/main.go

.PHONY: fmt
fmt: ## Run go fmt
	go fmt ./...

.PHONY: lint
lint: prereqs ## Lint go code
	@echo "### Linting go code"
	golangci-lint run ./...

.PHONY: test
test: ## Test go code
	@echo "### Testing go code"
	go test ./... -coverpkg=./... -coverprofile cover.out

##@ Images

.PHONY: image-build
image-build: ## Build MULTIARCH_TARGETS images
	trap 'exit' INT; \
	$(foreach target,$(MULTIARCH_TARGETS),$(call build_target,$(target)))

.PHONY: image-push
image-push: ## Push MULTIARCH_TARGETS images
	trap 'exit' INT; \
	$(foreach target,$(MULTIARCH_TARGETS),$(call push_target,$(target)))

.PHONY: manifest-build
manifest-build: ## Build MULTIARCH_TARGETS manifest
	@echo 'building manifest $(IMAGE)'
	DOCKER_BUILDKIT=1 $(OCI_BIN) rmi ${IMAGE} -f
	DOCKER_BUILDKIT=1 $(OCI_BIN) manifest create ${IMAGE} $(foreach target,$(MULTIARCH_TARGETS), --amend ${IMAGE}-$(target));

.PHONY: manifest-push
manifest-push: ## Push MULTIARCH_TARGETS manifest
	@echo 'publish manifest $(IMAGE)'
ifeq (${OCI_BIN}, docker)
	DOCKER_BUILDKIT=1 $(OCI_BIN) manifest push ${IMAGE};
else
	DOCKER_BUILDKIT=1 $(OCI_BIN) manifest push ${IMAGE} docker://${IMAGE};
endif
