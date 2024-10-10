# 快速开始

## 1. 定义服务协议

如: `example/proto/todo/todo.proto`

```proto
syntax = "proto3";

package todo;

option go_package = "github.com/Anderson-Lu/orion/proto_go/uit/todo";

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

## 2. 编译协议

```shell
cd example
make proto
```

## 3. 实现服务

如: `example/service/service.go`

```go
type Service struct {}

func (s *Service) Add(ctx context.Context, in *todo.AddReq) (*todo.AddRsp, error) {...}
func (s *Service) Remove(ctx context.Context, in *todo.RemoveReq) (*todo.RemoveRsp, error) {...}
func (s *Service) List(ctx context.Context, in *todo.ListReq) (*todo.ListRsp, error) {...}
func (s *Service) Add(ctx context.Context, in *todo.ModifyReq) (*todo.ModifyRsp, error) {...}
```

## 4. 运行服务

```go
package main

import (
  "log"

  "github.com/Anderson-Lu/orion/pkg/logger"
  "github.com/Anderson-Lu/orion/orpc"
  _ "github.com/Anderson-Lu/orion/orpc/build"

  "github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
  "github.com/Anderson-Lu/orion/example/orion_server/service"
)

func main() {

  c := &urpc.Config{
    Server:          &orpc.ServerConfig{Port: 8080, EnableGRPCGateway: true},
    PromtheusConfig: &orpc.PromtheusConfig{Enable: true, Port: 9092},
    FrameLogger:     &logger.LoggerConfig{Path: "../log/frame.log", LogLevel: "info"},
    AccessLogger:    &logger.LoggerConfig{Path: "../log/access.log"},
    ServiceLogger:   &logger.LoggerConfig{Path: "../log/service.log"},
    PanicLogger:     &logger.LoggerConfig{Path: "../log/panic.log"},
  }

  handler, _ := service.NewService(c)
  server, err := urpc.New(
    urpc.WithConfigFile("../config/config.toml"),
    urpc.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),
    urpc.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint),
    urpc.WithFlags(),
  )
  if err != nil {
    log.Fatal(err)
  }

  if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}
```

## 5. 构建服务

```shell
make build
```

## 6. 日志拆分

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

## 7. 提供HTTP服务

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