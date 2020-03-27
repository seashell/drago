SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

OS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH ?= amd64
CC ?= ""

ALL_TARGETS += linux_amd64 \
	linux_arm \
	linux_arm64 


GO_LDFLAGS ?= "-X github.com/seashell/drago/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"
CGO_ENABLED ?= 1

# Handle static builds
STATIC := $(or $(STATIC),$(S))
ifeq ($(STATIC),1)
	GO_LDFLAGS := "-linkmode external -extldflags -static"
endif


.PHONY: default
default: linux_amd64

.PHONY: linux_amd64
linux_amd64: ## Build linux_amd64 binary
	@echo "==> Building binary \"build/$@/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/$@/drago"	

.PHONY: linux_arm64
linux_arm64: ## Build linux_arm64 binary
	@echo "==> Building binary \"build/$@/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=${OS} GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/$@/drago"		

.PHONY: linux_arm
linux_arm: ## Build linux_arm binary
	@echo "==> Building binary \"build/$@/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/$@/drago"	

.PHONY: custom
custom: ## Build custom binary with the specified CC toolchain
	@echo "==> Building binary \"build/custom/${OS}_${ARCH}/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=${OS} GOARCH=${ARCH} CC=${CC} \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/custom/${OS}_${ARCH}/drago"	

.PHONY: ui
ui: ## Generate UI .go files
	go generate

.PHONY: clean
clean:
	go mod tidy
	rm -rf build
	rm -rf ui/build
	rm -rf ui/node_modules	

.PHONY: release
all: clean ui ${ALL_TARGETS} ## Clean build environment and build binaries for all targets with UI support

.PHONY: dev
dev: ui linux_amd64  ## Generate UI .go files and build linux_amd64 binary

# if vars not set specifially: try default to environment, else fixed value.
# strip to ensure spaces are removed in future editorial mistakes.
# tested to work consistently on popular Linux flavors and Mac.
ifeq ($(user),)
# USER retrieved from env, UID from shell.
HOST_USER ?= $(strip $(if $(USER),$(USER),nodummy))
HOST_UID ?= $(strip $(if $(shell id -u),$(shell id -u),4000))
else
# allow override by adding user= and/ or uid=  (lowercase!).
# uid= defaults to 0 if user= set (i.e. root).
HOST_USER = $(user)
HOST_UID = $(strip $(if $(uid),$(uid),0))
endif

.PHONY: docker  
docker: ## Build with Docker instead. e.g "make docker f=STATIC=1 c=all"
	@docker build --build-arg HOST_UID=${HOST_UID} --build-arg HOST_USER=${HOST_USER} -t drago_builder .
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} drago_builder /bin/sh -c ""${f}" make "${c}""

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""
	@echo "This host will build the following targets if 'make release' is invoked:"
	@echo $(ALL_TARGETS) | sed 's/^/    /'