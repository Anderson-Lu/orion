
# 自动注入构建版本

支持以Makefile方式打包二进制程序并动态注入框架版本等信息, UIT内置了`github.com/orion/urpcbuild`包,提供注入支持,当然这是可选的,或者按需自定义实现自己的build注入

```makefile
# makefile example
git_rev  = $(shell git rev-parse --short HEAD)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
app_name   = "example"

BuildVersion := $(git_branch)_$(git_rev)
BuildTime := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BuildCommit := $(shell git rev-parse --short HEAD)
BuildGoVersion := $(shell go version)
BuilderPkg := "github.com/orion/urpcbuild"

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
```

打印构建信息:

```shell
[root@VM /data/uit/example]# ./build/example -v
Git Branch   : dev_0_0_2_cbb23d5 
Git Commit   : cbb23d5 
Built Time   : 2024-09-23T12:21:35Z 
Go Version   : go version go1.22.1 linux/amd64 
Uit Version  : dev0.0.2 
```