#TODO: check if crosscompiling toolchains are present when running crossbuilds

# to enable arm builds: apt-get install -y gcc-arm-linux-gnueabihf libc6-dev-armhf-cross
# to enable arm64 builds: apt-get install -y gcc-aarch64-linux-gnu libc6-dev-arm64-cross


SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
THIS_OS := $(shell uname | tr '[:upper:]' '[:lower:]')

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_TEST_CMD = $(if $(shell which gotestsum),gotestsum --,go test)

GO_LDFLAGS := "-linkmode external -extldflags -static -X github.com/seashell/drago/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"

ALL_TARGETS = linux_amd64 \
	linux_arm \
	linux_arm64

OS ?= linux
ARCH ?= amd64

.PHONY: dev
dev:
	GOOS=${THIS_OS} GOARCH=${ARCH} go build -o build/${OS}-${ARCH}/drago main.go

.PHONY: ui
ui:
	go generate

build/linux_amd64/drago: $(SOURCE_FILES) ## Build Drago for linux/amd64
	@echo "==> Building $@ ..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64\
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "$@"

build/linux_arm/drago: $(SOURCE_FILES) ## Build Drago for linux/arm
	@echo "==> Building $@ ..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc-7 \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "$@"

build/linux_arm64/drago: $(SOURCE_FILES) ## Build Drago for linux/arm64
	@echo "==> Building $@ ..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc-7 \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "$@"	

.PHONY: release
release: clean ui $(foreach t,$(ALL_TARGETS),build/$(t)/drago) ## Build all release packages which can be built on this platform.

.PHONY: clean
clean:
	rm -rf build