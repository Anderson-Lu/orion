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

## 快速开始

```shell
git clone https://github.com/Anderson-Lu/uit.git

cd example/cmd

go run main.go
```

## 容器性能优化

在Golang的GPM模型中P的数量决定了并行的G的数量,而P的数量又是间接依赖M的数量,M的数量又由CPU的核心数量决定.在Golang中,通过`runtime.NumCPU()`获取宿主机的CPU核数,再通过`runtime.GOMAXPROCS(n)`设置GOMAXPROCS的值.当服务部署在容器中时,每个容器拿到的都是宿主机的核数,因此,当容器数量过多时,会导致产生过多的P,从而导致频繁的线程切换,最终导致服务性能下降.

因此,我们需要通过动态限制容器中GOMAXPROCS的值来避免上述问题.容器化是通过cgroup机制来限制容器能使用的cpu核心数的,因此,通过读取虚拟化为容器分配的cpu核数来为golang程序动态设置GOMAXPROCS的值. 这里引入`go.uber.org/automaxprocs`来实现容器内部署的性能优化.


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
