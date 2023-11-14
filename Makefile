SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
# gnu date format iso-8601 is parsable with Go RFC3339
BUILDDATE := $(shell date --iso-8601=seconds)
VERSION := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))

all: test cli

.PHONY: cli
cli:
	go build -ldflags "-X 'github.com/metal-stack/v.Version=$(VERSION)' \
								   -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
								   -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
								   -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'" \
	   -o bin/metal github.com/metal-stack-cloud/cli
	strip bin/metal

.PHONY: test
test:
	CGO_ENABLED=1 go test ./... -race -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out

.PHONY: golint
golint:
	golangci-lint run -p bugs -p unused -D protogetter
