**Orion**支持多种格式的配置文件, 如`json`,`yaml`和`toml`. 在初始化框架服务时指定即可:

```go
server, err := orpc.New(
  // uit.WithConfigFile("../config/config.toml"),
  // uit.WithConfigFile("../config/config.json"),
  // uit.WithConfigFile("../config/config.yaml"),
  // uit.WithConfig(&orpc.Config{...}),
  ...
)
```

配置项举例(yaml):

```yaml
# 服务配置
Server:
  # 服务GRPC端口
  Port: 8080
  # 是否同时兼容HTTP请求
  EnableGRPCGateway: true

# Promtheus相关的配置,会启动9092端口供相关的metrics采集
PromtheusConfig:
  Enable: true
  Port: 9092

# 框架本身的日志
FrameLogger:
  Path: '../log/frame.log'
  # 日志等级
  LogLevel: 'info'
  # 单文件最大值
  LogFileMaxSizeMB: 10
  # 最多保存的文件数
  LogFileMaxBackups: 10
  # 最大保存天数
  LogMaxAgeDays: 10
  # 是否开启gzip压缩
  LogCompress: false

# 服务接受流量的输入和输出日志
AccessLogger:
  Path: '../log/access.log'

# 服务发生panic时会捕获并记录
PanicLogger:
  Path: '../log/panic.log'

# 服务自身的业务日志
ServiceLogger:
  Path: '../log/service.log'

# 微服务注册配置,可选
# 需要在代码中指定对应的注册中心
# 如: orpc.WithRegistry(registry.RegisteyConsul)
Registry:
  # 当前服务的微服务名,会注册到对应的注册中心
  Service: "mine.namespace.demo"
  # 注册中心的IP地址
  IP: "127.0.0.1"
  # 注册中心的端口
  Port: 8500

# 接口限流配置
RateLimit:
 # 要限流的方法名
 - Key: "/todo.UitTodo/Add"
   # 令牌桶的容量
   Cap: 1
   # 令牌桶发放的速度
   TokensPerSecond: 1
```

更多配置模版:

[JSON格式配置](../example/uit_grpc_server/config/config.json) | [TOML格式配置](../example/uit_grpc_server/config/config.toml) | [YAML格式配置](../example/uit_grpc_server/config/config.yaml)

其中日志可选配置项:

```bash
Path # 文件路径
LogFileMaxSizeMB  # 单个文件最大限制(MB)
LogFileMaxBackups # 最多保存n个日志文件
LogMaxAgeDays # 最长保留天数
LogCompress  # 是否压缩(gzip)
LogLevel # 日志等级
```

**Orion**日志基于开源的`zaplog`和`lumberjack`组件进行二次组合封装.