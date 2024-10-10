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
Server:
  Port: 8080
  EnableGRPCGateway: true

PromtheusConfig:
  Enable: true
  Port: 9092

FrameLogger:
  Path: '../log/frame.log'
  LogLevel: 'info'

AccessLogger:
  Path: '../log/access.log'

PanicLogger:
  Path: '../log/panic.log'

ServiceLogger:
  Path: '../log/service.log'
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