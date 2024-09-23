## 什么是UIT
```
  __  __  ____  ____ 
 (  )(  )(_  _)(_  _)
  )(__)(  _)(_   )(  
 (______)(____) (__)  v0.0.1
```
UIT是一个基于GRPC的微服务框架，通过UIT可以快速构建同时支持GRPC以及HTTP的微服务，极大提升业务开发效率。

## 功能预览

- **协议兼容** 完全适配GRPC协议，同时提供HTTP访问兼容。
- **熔断限流** 支持接口层级以及用户自定义层级的限流策略，保证服务高可用。
- **日志集成** 集成zapLog/lumberjack日志，支持按照框架层、网关层、业务层分层输出，同时支持日志文件定期分割和清理清理等能力。
- **可观测性** 集成pprof metrics指标上报,接入Prometheus后可观测服务/接口等维度的各项健康指标(CPU/Memory/Routines/Fds等)。
- **链路追踪** 支持分布式链路追踪,帮助更快更方便排查和定位服务潜在问题。

## 性能优化项

- **自动设置服务最大进程** 自动设置同时使用的CPU核数,提升服务RPS.

## 快速开始

```shell
git clone https://github.com/Anderson-Lu/uit.git

cd example/cmd

go run main.go
```

## 自动注入构建版本

支持以Makefile方式打包二进制程序并动态注入框架版本等信息, UIT内置了`github.com/uit/pkg/xgrpc/build`包,提供注入支持,当然这是可选的,或者按需自定义实现自己的build注入

```makefile
# makefile example
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
```

打印构建信息:

```shell
[root@VM /data/uit/example]# ./build/example -v
Git Branch     : dev_0_0_2_cbb23d5 
Git Commit     : cbb23d5 
Built Time     : 2024-09-23T12:21:35Z 
Go Version     : go version go1.22.1 linux/amd64 
XGRPC Version  : dev0.0.2 
```

## 版本记录
