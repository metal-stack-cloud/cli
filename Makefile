GOOS := linux
GOARCH := amd64
BINARY := metal-$(GOOS)-$(GOARCH)
TAGS := -tags 'netgo'
SHA := $(shell git rev-parse --short=8 HEAD)
GITVERSION := $(shell git describe --long --all)
# gnu date format iso-8601 is parsable with Go RFC3339
BUILDDATE := $(shell date --iso-8601=seconds)
VERSION := $(or ${VERSION},$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD))

ifeq ($(GOOS),linux)
	LINKMODE := -linkmode external -extldflags '-static -s -w'
	TAGS := -tags 'osusergo netgo static_build'
endif

LINKMODE := $(LINKMODE) \
		 -X 'github.com/metal-stack/v.Version=$(VERSION)' \
		 -X 'github.com/metal-stack/v.Revision=$(GITVERSION)' \
		 -X 'github.com/metal-stack/v.GitSHA1=$(SHA)' \
		 -X 'github.com/metal-stack/v.BuildDate=$(BUILDDATE)'

all: test cli

.PHONY: cli
cli:
	go build \
		$(TAGS) \
		-ldflags \
		"$(LINKMODE)" \
		-o bin/$(BINARY) \
		github.com/metal-stack-cloud/cli
	md5sum bin/$(BINARY) > bin/$(BINARY).md5

.PHONY: test
test:
	CGO_ENABLED=1 go test ./... -race -coverprofile=coverage.out -covermode=atomic && go tool cover -func=coverage.out

.PHONY: golint
golint:
	golangci-lint run -p bugs -p unused -D protogetter
