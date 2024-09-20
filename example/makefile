git_rev    = $(shell git rev-parse --short HEAD)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
app_name   = "example"

BuildVersion := $(git_branch)_$(git_rev)
BuildTime := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BuildCommit := $(shell git rev-parse --short HEAD)
BuildGoVersion := $(shell go version)
BuilderPkg := "github.com/uit/pkg/xgrpc/build"

GOLDFLAGS =  -X '$(BuilderPkg).BuildVersion=$(BuildVersion)'
GOLDFLAGS += -X '$(BuilderPkg).BuildTime=$(BuildTime)'
GOLDFLAGS += -X '$(BuilderPkg).BuildCommit=$(BuildCommit)'
GOLDFLAGS += -X '$(BuilderPkg).BuildGoVersion=$(BuildGoVersion)'

.PHONY: build clean build-version

build:
	go build -ldflags "$(GOLDFLAGS)" -o build/$(app_name) cmd/main.go 
clean:
	rm build/*

build-version:
	./build/"$(app_name)" -v