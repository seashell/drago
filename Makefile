SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

THIS_OS := $(shell uname | cut -d- -f1)
THIS_ARCH := $(shell uname -m)

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_LDFLAGS ?= -X=github.com/seashell/drago/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)

CGO_ENABLED ?= 0

ifeq ($(CI),true)
	$(info Running in a CI environment, verbose mode is disabled)
else
	VERBOSE="true"
endif

# List of supported OS
SUPPORTED_OSES = Linux

# User defined flags
OS := $(or $(OS),$(O)) ## Define build target OS, e.g linux (coming soon)
ARCH := $(or $(ARCH),$(A)) ## Define build target architecture, e.g amd64 (coming soon)
STATIC := $(or $(STATIC),$(S)) ## If set to 1, build statically linked binary
DOCKER := $(or $(DOCKER),$(D)) ## If set to 1, run build within a Docker container

# ---- Handle static builds
ifeq ($(STATIC),1)
	GO_LDFLAGS := "${GO_LDFLAGS} -w -extldflags -static"
endif
	
# ---- In case of a Dockerized build, check if the builder image is available. 
ifeq ($(DOCKER),1)
	HOST_USER ?= $(strip $(if $(USER),$(USER),nodummy))
	HOST_UID ?= $(strip $(if $(shell id -u),$(shell id -u),4000))
	DOCKER_BUILDER_IMAGE_AVAILABLE := $(shell docker images --filter LABEL=com.drago.builder=true -q)
	BUILD_DOCKER_BUILDER_IMAGE_CMD := (docker build --label com.drago.builder=true --build-arg HOST_UID=${HOST_UID} --build-arg HOST_USER=${HOST_USER} -t drago_builder  . -f ./build/Dockerfile.builder)
endif

# =========== Targets ===========

ifeq (Linux,$(THIS_OS))
ALL_TARGETS = linux_amd64
endif

default: help

# ====> Current platform
.PHONY: dev
dev: GOOS=$(shell go env GOOS)
dev: GOARCH=$(shell go env GOARCH)
dev: DEV_TARGET=$(GOOS)_$(GOARCH)
dev: ## Build for the current platform
	@rm -rf $(PROJECT_ROOT)/bin
	@$(MAKE) --no-print-directory $(DEV_TARGET)

# ====> Container
.PHONY: container
container: ## Build container with the Drago binary inside
	@$(MAKE) ui dev STATIC=1
	@echo "==> Building container image "drago:latest" ..."
	@docker build -t drago:latest . -f ./build/Dockerfile.linux_amd64

# ====> All
.PHONY: all
all: clean ui $(foreach t,$(ALL_TARGETS),$(t)) ## Build all targets supported by this platform
	@echo "==> Results:"
	@tree --dirsfirst $(PROJECT_ROOT)/bin

# ====> Tidy
.PHONY: tidy
tidy:
	@echo "--> Tidying up Drago modules"
	@go mod tidy

# ====> Linux AMD 64
.PHONY: linux_amd64
linux_amd64: CMD='CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 \
								go build \
								-trimpath \
								-ldflags ""$(GO_LDFLAGS)"" \
								-o "bin/$@/drago"'
linux_amd64: $(SOURCE_FILES) ## Build Drago for linux/amd64
	@echo "==> Building $@..."
	@echo "==> COMMAND $(CMD)..."
ifeq ($(DOCKER),1)
ifeq ($(DOCKER_BUILDER_IMAGE_AVAILABLE),)
	@echo "==> Building Docker builder image..."
	@$(call BUILD_DOCKER_BUILDER_IMAGE_CMD)
endif	
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} drago_builder \
	/bin/sh -c ${CMD}
else
	@eval ${CMD}
endif

# ====> Linux ARM 64
.PHONY: linux/arm64
linux/arm64: ## Build Drago for linux/arm64 (coming soon)
	@echo "==> Coming soon..."

# ====> Linux ARM
.PHONY: linux/arm
linux/arm: ## Build drago for linux/arm (coming soon)
	@echo "==> Coming soon..."

# ====> Web UI
.PHONY: ui
ui: CMD="go generate"
ui: ## Build Web UI
	@echo "==> Building Web UI..."
ifeq ($(DOCKER),1)
ifeq ($(DOCKER_BUILDER_IMAGE_AVAILABLE),)
	@echo "==> Building Docker builder image..."
	@$(call BUILD_DOCKER_BUILDER_IMAGE_CMD)
endif	
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} drago_builder \
	/bin/sh -c ${CMD}
else
	@eval ${CMD}
endif

.PHONY: clean
clean: ## Remove build artifacts
	@echo "==> Cleaning build artifacts..."
	@rm -rf "$(PROJECT_ROOT)/bin/"
	@rm -rf "$(PROJECT_ROOT)/ui/build/*"
	@rm -rf "$(PROJECT_ROOT)/ui/node_modules/"

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
EG_FORMAT="    \033[36m%s\033[0m %s\n"

.PHONY: help
help: ## Display usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""
	@echo "This host will build the following targets if 'make release' is invoked:"
	@echo $(ALL_TARGETS) | sed 's/^/    /'
	@echo ""
	@echo "Valid flags:"
	@grep -E '^[^ ]+ :=.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = " :=.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make dev STATIC=1"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make ui dev DOCKER=1"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make container DOCKER=1"
