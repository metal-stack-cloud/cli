
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
BUILDDATE := $(shell date -Iseconds)
VERSION := $(or ${VERSION},devel)

all: cli

.PHONY: cli
cli: test
	go build -tags netgo -ldflags "-X 'github.com/metal-stack/v.Version=$(VERSION)' \
								   -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
								   -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
								   -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'" \
	   -o bin/cli github.com/metal-stack-cloud/cli
	strip bin/cli

.PHONY: test
test:
	CGO_ENABLED=1 go test ./... -race -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out

.PHONY: golint
golint:
	golangci-lint run -p bugs -p unused