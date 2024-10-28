可观测性大致可以分为以下三种类型:

- 日志(Log)
- 指标(Metrics)
- 链路追踪(Trace)

在传统的微服务可观测性实践上, 往往将上述三种类型的数据各自单独出来设计一套工作流, 比如

- 日志记录有ELK方案, ELK + Filebeat, ELK + Logstash等等.
- 指标方案有Prothemus等等.
- 链路追踪方案有Opentracing等等.

Orion则借助OpenTelemetry实现工作流的统一.