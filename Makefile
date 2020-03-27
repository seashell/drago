SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))


GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_TEST_CMD = $(if $(shell which gotestsum),gotestsum --,go test)

OS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH ?= amd64
CC ?= ""


GO_LDFLAGS ?= "-X github.com/seashell/drago/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"
CGO_ENABLED ?= 1

# Handle static builds
STATIC := $(or $(STATIC),$(S))
ifeq ($(STATIC),1)
	GO_LDFLAGS := "-linkmode external -extldflags -static"
	CGO_ENABLED := 0
endif


.PHONY: default
default: amd64

.PHONY: amd64
amd64:
	@echo "==> Building binary \"build/${OS}_${ARCH}/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=${OS} GOARCH=$@ \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/${OS}_${ARCH}/drago" \
		main.go

.PHONY: arm64
arm64:
	@echo "==> Building binary \"build/${OS}_$@/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=${OS} GOARCH=$@ CC=aarch64-linux-gnu-gcc \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/${OS}_$@/drago"	

.PHONY: arm
arm:
	@echo "==> Building binary \"build/${OS}_$@/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=${OS} GOARCH=$@ CC=arm-linux-gnueabihf-gcc \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/${OS}_$@/drago"	

.PHONY: custom
custom:
	@echo "==> Building binary \"build/custom/${OS}_${ARCH}/drago\" ..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=${OS} GOARCH=${ARCH} CC=${CC} \
		go build \
		-ldflags $(GO_LDFLAGS) \
		-o "build/${OS}_${ARCH}/drago"	

.PHONY: ui
ui:
	go generate

.PHONY: clean
clean:
	go mod tidy
	rm -rf build
	rm -rf ui/build
	rm -rf ui/node_modules	

.PHONY: all
all: clean ui amd64 arm64 arm

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
docker:
	@docker build --build-arg HOST_UID=${HOST_UID} --build-arg HOST_USER=${HOST_USER} -t drago_builder .
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} drago_builder /bin/sh -c ""${f}" make "${c}""
