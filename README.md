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

## 性能优化项

- **自动设置服务最大进程** 自动设置同时使用的CPU核数,提升服务RPS.

## 快速开始

```shell
git clone https://github.com/Anderson-Lu/uit.git

cd example/cmd

go run main.go
```

## 版本记录
