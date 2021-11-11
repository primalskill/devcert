SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

BUILDDATE := $(shell date +'%Y-%m-%d')
BUILDTIME := $(shell date +'%H:%M:%S')

CURPATH := ${shell pwd}

APP := devcert
#APP_INPUT := ${CURPATH}/main.go
APP_INPUT := .
APP_OUTPUT := /usr/local/bin/${APP}

RELEASE_OUTPUT := ./.bin/${APP}
RELEASE_WIN_AMD64_EXT := ${RELEASE_OUTPUT}_win_amd64.exe
RELEASE_DARWIN_AMD64_EXT := ${RELEASE_OUTPUT}_darwin_amd64
RELEASE_DARWIN_ARM64_EXT := ${RELEASE_OUTPUT}_darwin_arm64
RELEASE_LINUX_AMD64_EXT := ${RELEASE_OUTPUT}_linux_amd64
RELEASE_LINUX_ARM64_EXT := ${RELEASE_OUTPUT}_linux_arm64

.PHONY: clean format release

update-deps:
	go mod tidy
	go get -u ./...

clean:
	go clean ${APP_INPUT}

format:
	gofmt -w ${CURPATH}

build:
	go build \
	-ldflags "-X main.BuildDate=$(BUILDDATE) -X main.BuildTime=$(BUILDTIME)" \
	-o $(APP_OUTPUT) \
	$(APP_INPUT)

release-win-amd64:
	env GOOS=windows GOARCH=amd64 \
	go build \
	-ldflags "-X main.BuildDate=$(BUILDDATE) -X main.BuildTime=$(BUILDTIME)" \
	-o $(RELEASE_WIN_AMD64_EXT) \
	$(APP_INPUT)

release-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "-X main.BuildDate=$(BUILDDATE) -X main.BuildTime=$(BUILDTIME)" \
	-o $(RELEASE_DARWIN_AMD64_EXT) \
	$(APP_INPUT)


release-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 \
	go build \
	-ldflags "-X main.BuildDate=$(BUILDDATE) -X main.BuildTime=$(BUILDTIME)" \
	-o $(RELEASE_DARWIN_ARM64_EXT) \
	$(APP_INPUT)

release-linux-amd64:
	env GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags "-X main.BuildDate=$(BUILDDATE) -X main.BuildTime=$(BUILDTIME)" \
	-o $(RELEASE_LINUX_AMD64_EXT) \
	$(APP_INPUT)

release-linux-arm64:
	env GOOS=linux GOARCH=arm64 \
	go build \
	-ldflags "-X main.BuildDate=$(BUILDDATE) -X main.BuildTime=$(BUILDTIME)" \
	-o $(RELEASE_LINUX_ARM64_EXT) \
	$(APP_INPUT)


compile-releases: release-win-amd64 release-darwin-amd64 release-darwin-arm64 release-linux-amd64 release-linux-arm64
	echo "Compiled releases"

exec: clean format build
	${APP_OUTPUT}