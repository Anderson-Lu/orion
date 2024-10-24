# 1. 创建GRPC客户端

```go
import "github.com/Anderson-Lu/orion/orpc/client"

cli, err := client.New(resolver)
```

# 2. Orion请求选项(OrionClientInvokeOption)

Orion内置了很多请求选项, 业务按需调用即可:

```go
import "github.com/Anderson-Lu/orion/orpc/client/options"
```

|选项|功能|
|:-|:-|
|`options.WithJson()`|表示请求体的序列化协议为json,开启此项需要服务端支持(Orion服务端默认开启),否则默认使用proto协议进行序列化|
|`options.WithCircuitBreak()`|表示该请求接入熔断检测,需要前置注册对应的熔断器(cli.RegisterCircuitBreakRule())|
|`options.WithService()`|表示通过微服务名进行寻址,需要前置指定对应的resolver为非直连方式|
|`options.WithMethod()`|表示要调用的方法名,指pb生成后的完整PathDesc|
|`options.WithDirectAddress()`|表示本次请求直接通过指定的地址访问,如一些调试环境,需要指定特定IP/端口时|
|`options.WithHeaders()`|将kv添加到OutGoingMetadata里面|

# 3. 序列化

- 客户端在调用的时候支持编码器`client.WithJson()`

```go
import "github.com/Anderson-Lu/orion/orpc/client"

// ...

if err := cli.Invoke(ctx, req, rsp, client.WithJson()); err != nil {
  // TODO
    // add your logic
}
```

- 服务端注册对应的编码器(Orion内置默认开启):

```go
import _ "github.com/Anderson-Lu/orion/orpc/codec"
```

# 4. 配置寻址/服务发现(OrionResolver)

在Orion中,服务发现统一封装在不同的**Resolver**中, 用于服务的发现:

```go
type IResolver interface {
	Name() string
	Select(serviceName string) (*grpc.ClientConn, error)
}
```

因此,业务可以实现自己的微服务服务发现逻辑,从而保证Orion框架的灵活性和扩展性. 目前,内置了直连模式(DirectResolver)和Consul模式(ConsulResover)

## 4.1 直连IP模式(DirectResolver)

```go
// 127.0.0.1:8080是GRPC服务端的服务端口和地址,resolver.NewDirectResolver()表示通过直连方式实现服务发现
cli, err := client.New(resolver.NewDirectResolver("127.0.0.1:8080"))

// 请求选项
opts := []options.OrionClientInvokeOption{
  options.WithMethod("/todo.UitTodo/Add"),
}

// 执行请求
err := cli.Invoke(context.Background(), req, rsp, opts...)
```

## 4.2 Consul服务发现模式(ConsulResover)

```go
// 127.0.0.1:8500是consul agent的IP和端口, resolver.NewConsulResovler表示通过consul方式实现动态服务发现
cli, err := client.New(resolver.NewConsulResovler("127.0.0.1:8500"))

// 请求选项
opts := []options.OrionClientInvokeOption{
  // 表示要进行寻址的微服务名(微服务ID)
  options.WithService("mine.namespace.demo"),
  options.WithMethod("/todo.UitTodo/Add"),
}

// 执行请求
err := cli.Invoke(context.Background(), req, rsp, opts...)
```

# 5. 配置熔断器

依赖[Orion熔断组件](./doc_circuit_breaker.md), Orion框架内部集成熔断策略. 只需要在创建客户端时指定对应的选项即可.

```go
// 直连模式
rsv := resolver.NewDirectResolver("127.0.0.1:8080")
cli, err := client.New(rsv)

// 创建客户端并注册对应的熔断策略
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

# 6. 配置负载均衡(Balancer)

Orion支持在主调侧设置需要服务发现的负载均衡器,注意,这里不用传统意义上的连接池, 原因如下:

- GRPC连接本身是基于http2的,本身有多路复用等特性,client本身也有有重试、断线重连的能力,因此连接池的作用其实不大.
- 如果采用连接池,在一些特定场景之下,会影响负载均衡的实现.

Orion客户端支持三种Balancer:

- 随机`random`,在已创建的连接中随机选择一个.
- 哈希`hash`,通过指定字段固定映射到一个连接.
- 默认`default`, 轮转的方式(Round Robin), 每次调用依次选择.

当前Orion客户端默认使用随机LB策略(balancer.DefaultBalancer)来充当服务发现后的连接选择决策.如果要通过其他LB策略,则需要在创建对应resolver对象时指定,如`ConsulResovler`

```go
b := balancer.NewDefaultBalancer() // 创建一个默认的随机LB
c := resolver.NewConsulResovler("127.0.0.1:8500",b) // 创建一个Consul服务发现客户端,并指定LB策略为随机
```

以下是Orion内置的LB策略:

|Balancer|均衡算法|
|:-|:-|
|balancer.RoundRobinBalancer|轮询,默认模式|
|balancer.Crc32HashBalancer|一致性哈希,需要注意IP漂移问题|

# 7. 超时控制(Timeout)

Orion不再单独提供超时相关的Option,直接通过context来实现超时控制.

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

err := cli.Invoke(ctx, req, rsp, opts...)
```

超时则按照GRPC的统一错误码返回对应的错误(code=4):

```
rpc error: code = DeadlineExceeded desc = context deadline exceeded
```
