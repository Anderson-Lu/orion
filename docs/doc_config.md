`Orion`支持多种格式的配置文件, 如`json`,`yaml`和`toml`. 在初始化框架服务时指定即可:

```go
server, err := uit.New(
  // uit.WithConfigFile("../config/config.toml"),
  // uit.WithConfigFile("../config/config.json"),
  // uit.WithConfigFile("../config/config.yaml"),
  // uit.WithConfig(&uit.Config{...}),
  ...
)
```

配置项举例:

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

其中日志可选配置项:

```go
type LoggerConfig struct {
  // 文件路径
	Path              string `yaml:"Path" json:"Path" toml:"Path"`
  // 单个文件最大限制(MB)
	LogFileMaxSizeMB  int    `yaml:"LogFileMaxSizeMB" json:"LogFileMaxSizeMB" toml:"LogFileMaxSizeMB"`
	// 最多保存n个日志文件
  LogFileMaxBackups int    `yaml:"LogFileMaxBackups" json:"LogFileMaxBackups" toml:"LogFileMaxBackups"`
	// 最长保留天数
  LogMaxAgeDays     int    `yaml:"LogMaxAgeDays" json:"LogMaxAgeDays" toml:"LogMaxAgeDays"`
	// 是否压缩(gzip)
  LogCompress       bool   `yaml:"LogCompress" json:"LogCompress" toml:"LogCompress"`
	// 日志等级
  LogLevel          string `yaml:"LogLevel" json:"LogLevel" toml:"LogLevel"`
}
```