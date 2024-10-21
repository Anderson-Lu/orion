## 什么是Orion

**Orion**是一个基于**GRPC**的微服务框架，通过**Orion**可以快速构建GRPC微服务.

## 功能预览

- **协议兼容** 完全适配GRPC协议，同时提供HTTP访问兼容。
- **熔断限流** 支持接口层级以及用户自定义层级的限流策略，保证服务高可用。
- **日志集成** 支持按照框架层、业务层分层输出，同时支持日志文件定期分割和清理清理等能力。
- **可观测性** 集成pprof metrics指标上报,接入Prometheus后可观测服务/接口等维度的各项健康指标(CPU/Memory/Routines/Fds等)。
- **链路追踪** 支持分布式链路追踪,帮助更快更方便排查和定位服务潜在问题。
- **工具集成** 集成常用的工具集

## 快速开始

- [快速开始](./docs/doc_get_start.md)
- [Orion设计哲学](./docs/doc_design.md)
- [Orion配置说明](./docs/doc_config.md) 
- [Orion服务注册&发现](./docs/doc_discovery.md)
- [Orion-Cli脚手架](./docs/doc_cli.md) 
- [Orion客户端](./docs/doc_orion_client.md) 
- [Orion错误码](./orpc/codes/codes.go)

## Orion优雅编程

- [熔断器(OrionCircuitBreaker)](./docs/doc_circuit_breaker.md)
- [限流器(OrionRateLimiter)](./docs/doc_get_start.md)
- [分发器(OrionDispatcher)](./docs/doc_get_start.md)
- [容器优化](./docs/doc_docker.md)
- [本地缓存TickerCache](./docs/doc_ticker_cache.md)
- [自动注入构建版本](./docs/doc_build_tool.md)
