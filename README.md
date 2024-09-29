## 什么是UIT
```
  __  __  ____  ____ 
 (  )(  )(_  _)(_  _)
  )(__)(  _)(_   )(  
 (______)(____) (__)  v0.0.2
```
**UIT**是一个基于**GRPC**的微服务框架，通过**UIT**可以快速构建GRPC微服务.

- 避免重复造轮子,提升业务开发效率
- 优雅编程

## 功能预览

- **协议兼容** 完全适配GRPC协议，同时提供HTTP访问兼容。
- **熔断限流** 支持接口层级以及用户自定义层级的限流策略，保证服务高可用。
- **日志集成** 支持按照框架层、业务层分层输出，同时支持日志文件定期分割和清理清理等能力。
- **可观测性** 集成pprof metrics指标上报,接入Prometheus后可观测服务/接口等维度的各项健康指标(CPU/Memory/Routines/Fds等)。
- **链路追踪** 支持分布式链路追踪,帮助更快更方便排查和定位服务潜在问题。

## 快速开始

#### 1. 定义服务协议

如: `example/proto/todo/todo.proto`

```proto
syntax = "proto3";

package todo;

option go_package = "github.com/uit/proto_go/uit/todo";

// TodoStatus 状态
enum TodoStatus {
  TodoStatusNone = 0;
  TodoStatusStarted = 1;
  TodoStatusFinished = 2;
}

// TodoItem 条目
message TodoItem {
  string id = 1;
  string title = 2;
  string desc = 3;
  repeated string tags = 4;
}

message AddReq {}

message AddRsp {}

message ListReq {}

message ListRsp {
  repeated TodoItem items = 1;
}

message RemoveReq {}

message RemoveRsp {}

message ModifyReq {}

message ModifyRsp{}

service UitTodo {
  rpc Add(AddReq) returns (AddRsp);
  rpc Remove(RemoveReq) returns (RemoveRsp);
  rpc List(ListReq) returns (ListRsp);
  rpc Modify(ModifyReq) returns (ModifyRsp);
}
```

#### 2. 编译协议

```shell
cd example
make proto
```

#### 3. 实现服务

如: `example/service/service.go`

```go
type Service struct {}

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {...}
func (s *Service) Remove(ctx context.Context, in *todo.RemoveReq) (*todo.RemoveRsp, error) {...}
func (s *Service) List(ctx context.Context, in *todo.ListReq) (*todo.ListRsp, error) {...}
func (s *Service) Add(ctx context.Context, in *todo.ModifyReq) (*todo.ModifyRsp, error) {...}
```

#### 4. 运行服务

```go
package main

import (
	"log"

	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/uit"
	_ "github.com/uit/pkg/uit/build"

	"github.com/uit/example/proto_go/proto/todo"
	"github.com/uit/example/service"
)

func main() {

	c := &uit.Config{
		Server:          &uit.ServerConfig{Port: 8080, EnableGRPCGateway: true},
		PromtheusConfig: &uit.PromtheusConfig{Enable: true, Port: 9092},
		FrameLogger:     &logger.LoggerConfig{Path: []string{"..", "log", "frame.log"}, LogLevel: "info"},
		AccessLogger:    &logger.LoggerConfig{Path: []string{"..", "log", "access.log"}},
		ServiceLogger:   &logger.LoggerConfig{Path: []string{"..", "log", "service.log"}},
		PanicLogger:     &logger.LoggerConfig{Path: []string{"..", "log", "panic.log"}},
	}

	handler, _ := service.NewService(c)
	server, err := uit.New(c,
		uit.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),
		uit.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint),
		uit.WithFlags(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
```

##### 5. 构建服务

```shell
make build
```

## 服务发现和注册

微服务框架本身的服务注册与发现理论上可以从框架本身摘除,就目前的主流模式(云原生)而言,微服务的注册和发现更倾向于将其沉淀到sidecar上,比如`service mesh`方式的实现,并以此可以实现更多底层需求(流量控制,熔断限流等,模调监控),因此,对于不同企业,其实现方式也不尽相同,因此,uit框架本身不做固化的服务注册和发现逻辑,开发者可以自己的需要灵活实现.

当然uit框架在未来也可能会提供注册和发现的相关组件.

## 提供HTTP服务

使用`grpc-gateway`插件为服务提供http服务, 支持在同一个端口同时支持GRPC和HTTP协议,配置项:

```go
server, err := uit.New(c,
  uit.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),          // 支持 grpc 协议
  uit.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint), // 支持 http 协议
  ...
)
```

之后可以使用http访问grpc服务:

```shell
curl -XPOST 'http://127.0.0.1:8080/todo.UitTodo/Add' -H 'Content-Type:application/grpc' -d '{"item":{"title":"title","desc":"desc","tags":["1","2"]}}'
```

使用grpctool访问grpc服务:

```shell
grpcurl -plaintext 127.0.0.1:8080 list
> todo.UitTodo

grpcurl -plaintext 127.0.0.1:8080 list todo.UitTodo
> todo.UitTodo.Add
> todo.UitTodo.List
> todo.UitTodo.Modify
> todo.UitTodo.Remove
```

## 日志拆分

uit默认支持以下4种日志,分别为:

- `frame`: 框架本身的一些日志
- `access`: 流量访问的一些日志,比如对应请求和回包的日志
- `panic`: 程序发生panic的一些捕获日志
- `service`: 业务逻辑自身的日志

代码配置方式:

```
c := &uit.Config{
  ...
  FrameLogger:   &logger.LoggerConfig{Path: []string{"..", "log", "frame.log"}, LogLevel: "info"},
  AccessLogger:  &logger.LoggerConfig{Path: []string{"..", "log", "access.log"}},
  ServiceLogger: &logger.LoggerConfig{Path: []string{"..", "log", "service.log"}},
  PanicLogger:   &logger.LoggerConfig{Path: []string{"..", "log", "panic.log"}},
  ...
}
```

**access.log** 举例:

```log
{"level":"info","time":"2024-01-29T10:58:42.613+0800","caller":"interceptors/accesslog.go:27","message":"[succ]","method":"/todo.UitTodo/Add","requestId":"","clientIP":"","req":"item:{title:\"title\"  desc:\"desc\"  tags:\"1\"  tags:\"2\"}","rsp":"","cost":0}
```

**frame.log** 举例:

```log
{"level":"info","time":"2024-01-29T10:58:39.841+0800","caller":"uit/server.go:175","message":"[Server] gRPC server started succ","port":8080}
```

## 容器性能优化

在Golang的GPM模型中P的数量决定了并行的G的数量,而P的数量又是间接依赖M的数量,M的数量又由CPU的核心数量决定.在Golang中,通过`runtime.NumCPU()`获取宿主机的CPU核数,再通过`runtime.GOMAXPROCS(n)`设置GOMAXPROCS的值.当服务部署在容器中时,每个容器拿到的都是宿主机的核数,因此,当容器数量过多时,会导致产生过多的P,从而导致频繁的线程切换,最终导致服务性能下降.

因此,我们需要通过动态限制容器中GOMAXPROCS的值来避免上述问题.容器化是通过cgroup机制来限制容器能使用的cpu核心数的,因此,通过读取虚拟化为容器分配的cpu核数来为golang程序动态设置GOMAXPROCS的值. 这里引入`go.uber.org/automaxprocs`来实现容器内部署的性能优化.


## 自动注入构建版本

支持以Makefile方式打包二进制程序并动态注入框架版本等信息, UIT内置了`github.com/uit/pkg/uit/build`包,提供注入支持,当然这是可选的,或者按需自定义实现自己的build注入

```makefile
# makefile example
git_rev  = $(shell git rev-parse --short HEAD)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
app_name   = "example"

BuildVersion := $(git_branch)_$(git_rev)
BuildTime := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BuildCommit := $(shell git rev-parse --short HEAD)
BuildGoVersion := $(shell go version)
BuilderPkg := "github.com/uit/pkg/uit/build"

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

## 错误码

UIT框架内置了一些错误码,原则上业务使用的错误码应当与框架的错误码区分开,详细信息可以参照`uit/pkg/uit/codes/codes.go`.
