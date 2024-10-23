# 创建GRPC客户端

```go
import "github.com/Anderson-Lu/orion/orpc/client"

cli, err := client.New(&client.OrionClientConfig{
  Host:               "127.0.0.1:8080",
  DailTimeout:        3000,
  ConnectionNum:      1,
  ConnectionBalancer: "random",
})
```

配置项:

```go
type OrionClientConfig struct {
  Host               string // server
  DailTimeout        int64  // milliseconds
  ConnectionNum      int64  // if set to zero, ConnectionNum = 1 will be set.
  ConnectionBalancer string // can be set to "random" or "hash", if not set, random balancer will be set by default.
}
```

# 使用Json协议发起grpc请求

- 客户端在调用的时候支持编码器`client.WithJson()`

```go
import "github.com/Anderson-Lu/orion/orpc/client"

// ...

if err := cli.Invoke(ctx, "your method desc", req, rsp, client.WithJson()); err != nil {
  // TODO
    // add your logic
}
```

- 服务端注册对应的编码器(Orion内置默认开启):

```go
import _ "github.com/Anderson-Lu/orion/orpc/codec"
```

# 主调侧GRPC长链接Balancer

Orion支持在主调侧设置需要创建的长链接(拨号连接)数量,注意,这里不用传统意义上的连接池, 原因如下:

- GRPC连接本身是基于http2的,本身有多路复用等特性,client本身也有有重试、断线重连的能力,因此连接池的作用其实不大.
- 如果采用连接池,在一些特定场景之下,会影响负载均衡的实现.

因此,Orion只对连接对象做简单的LB设置, 比如可以创建n个长链接对象,并通过不同均衡器(`ConnectionBalancer选项`)使用不同的连接进行调用.

Orion客户端支持三种Balancer:

- 随机`random`,在已创建的连接中随机选择一个.
- 哈希`hash`,通过指定字段固定映射到一个连接.
- 默认`default`, 轮转的方式, 每次调用依次选择.

其中, 如果选择了哈希模式, 则在调用过程中需要指定需要进行哈希计算的数据:

```go
// 需要指定进行hash的字段: client.WithHash("uid")
if err := cli.Invoke(ctx, "/todo.UitTodo/Add", req, rsp, client.WithHash("uid")); err != nil {
  panic(err)
}
```

注意, 当前的**Balancer**只吃对**本地**GRPC长链接对象的均衡分发策略, 非被调端的负载均衡.一般情况下,只需要维持一个长链接即可满足普通的业务需求.

# 主调侧熔断器

依赖[熔断组件](./doc_circuit_breaker.md), Orion框架内部集成熔断策略. 只需要在创建客户端时指定对应的选项即可.

```go
// 创建客户端并注册对应的熔断策略
cli, err := client.New(&client.OrionClientConfig{Host: "127.0.0.1:8080", DailTimeout: 10000})
cli.RegisterCircuitBreakRule(&circuit_break.RuleConfig{
  Name:             "/todo.UitTodo/Add",
  Window:           &circuit_break.WindowConfig{Duration: 1000, Buckets: 10},
  OpenDuration:     1000,
  HalfOpenDuration: 1000,
  HaflOpenPassRate: 0,
  RuleExpression:   "req_count > 0 && succ_rate < 100.0",
})

// Invoke选项指定配置熔断
opts := []client.OrionClientInvokeOption{
  client.WithJson(),
  client.WithCircuitBreak(), // 开启熔断检测
}

// 调用
err := cli.Invoke(context.Background(), "/todo.UitTodo/Add", req, rsp, opts...)

// 判断是否被熔断了(3001)
isCircuitError := codes.GetCodeFromError(err) == codes.ErrCodeCircuitBreak
```