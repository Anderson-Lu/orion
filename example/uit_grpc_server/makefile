git_rev    = $(shell git rev-parse --short HEAD)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
app_name   = "example"
proto_path = "./proto/todo"

BuildVersion := $(git_branch)_$(git_rev)
BuildTime := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BuildCommit := $(shell git rev-parse --short HEAD)
BuildGoVersion := $(shell go version)
BuilderPkg := "github.com/uit/pkg/uit/build"

GOLDFLAGS =  -X '$(BuilderPkg).BuildVersion=$(BuildVersion)'
GOLDFLAGS += -X '$(BuilderPkg).BuildTime=$(BuildTime)'
GOLDFLAGS += -X '$(BuilderPkg).BuildCommit=$(BuildCommit)'
GOLDFLAGS += -X '$(BuilderPkg).BuildGoVersion=$(BuildGoVersion)'

.PHONY: build clean build-version proto

build:
	go build -ldflags "$(GOLDFLAGS)" -o build/$(app_name) cmd/main.go 
clean:
	rm build/*

build-version:
	./build/"$(app_name)" -v

proto:
	for fsName in `ls "$(proto_path)"`; do echo "protoc >> ""$(proto_path)/$$fsName"; protoc --go_out=proto_go --go_opt=paths=source_relative --go-grpc_out=proto_go --go-grpc_opt=paths=source_relative --grpc-gateway_out=proto_go --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true $(proto_path)/$$fsName; done;
