# to enable arm builds: apt-get install -y gcc-arm-linux-gnueabihf libc6-dev-armhf-cross
# to enable arm64 builds: apt-get install -y gcc-aarch64-linux-gnu libc6-dev-arm64-cross


SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
THIS_OS := $(shell uname | tr '[:upper:]' '[:lower:]')

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_TEST_CMD = $(if $(shell which gotestsum),gotestsum --,go test)

GO_LDFLAGS ?= "-X github.com/seashell/drago/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"
CGO_ENABLED ?= 1

ALL_TARGETS = linux_amd64 \
	linux_arm \
	linux_arm64

OS ?= linux
ARCH ?= amd64

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

# Handle static builds
STATIC := $(or $(STATIC),$(S))
ifeq ($(STATIC),1)
	GO_LDFLAGS := "-linkmode external -extldflags -static ${GO_LDFLAGS}"
endif


.PHONY: dev
dev:
	GOOS=${THIS_OS} GOARCH=${ARCH} go build -o build/${OS}-${ARCH}/drago main.go

.PHONY: ui
ui:
	go generate

build/linux_amd64/drago: $(SOURCE_FILES) ## Build Drago for linux/amd64
	@echo "==> Building $@ ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64\
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "$@"

build/linux_arm/drago: $(SOURCE_FILES) ## Build Drago for linux/arm
	@echo "==> Building $@ ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "$@"

build/linux_arm64/drago: $(SOURCE_FILES) ## Build Drago for linux/arm64
	@echo "==> Building $@ ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "$@"	

.PHONY: clean
clean:
	rm -rf build
	rm -rf ui/build
	rm -rf ui/node_modules

.PHONY: _release
_release: clean ui $(foreach t,$(ALL_TARGETS),build/$(t)/drago) ## Build all release packages which can be built on this platform.

.PHONY: release
release:
ifeq ($(STATIC), 1)
	@echo "==> Static build..."
	@docker build --build-arg HOST_UID=${HOST_UID} --build-arg HOST_USER=${HOST_USER} -t static .
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} static /bin/sh -c "CGO_ENABLED=0 make _release"
else
	@echo "==> Regular build..."
	make _release
endif
	
