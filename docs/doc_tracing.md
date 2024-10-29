可观测性大致可以分为以下三种类型:

- 日志(Log)
- 指标(Metrics)
- 链路追踪(Trace)

在传统的微服务可观测性实践上, 往往将上述三种类型的数据各自单独出来设计一套工作流, 比如

- 日志记录有ELK方案, ELK + Filebeat, ELK + Logstash等等.
- 指标方案有Prothemus等等.
- 链路追踪方案有Opentracing等等.

Orion则借助OpenTelemetry实现工作流的统一.

# Metadata

|key|description|
|:-|:-|
|traceID|traceID|
|spanID|spanID|
|operationName|链路节点名|
|startTime|请求时间|
|duration|耗时|

# Tags

|key|description|
|:-|:-|
|deployment.environment|环境(production/test)|
|deployment.namespace|命名空间,按业务分类|
|service.instance.id|主调微服务名|
|service.instance.ip|主调IP|
|callee.service|被调微服务名|
|otel.status_description|错误信息|
|span.code|错误码|

Orion生成的Tracing数据举例如下:

```shell
{
  "data": [
    {
      "traceID": "43337abeaf760fe1ef2a02704fc7c12c",
      "spans": [
        {
          "traceID": "43337abeaf760fe1ef2a02704fc7c12c",
          "spanID": "942af23b7b7a9ba1",
          "operationName": "precheck",
          "references": [
            
          ],
          "startTime": 1730166469127984,
          "duration": 1000045,
          "tags": [
            {
              "key": "otel.library.name",
              "type": "string",
              "value": "主调微服务"
            },
            {
              "key": "callee.service",
              "type": "string",
              "value": "被调1"
            },
            {
              "key": "span.kind",
              "type": "string",
              "value": "internal"
            },
            {
              "key": "otel.status_code",
              "type": "string",
              "value": "ERROR"
            },
            {
              "key": "error",
              "type": "bool",
              "value": true
            },
            {
              "key": "otel.status_description",
              "type": "string",
              "value": "error occuro, 遇到错误了"
            },
            {
              "key": "internal.span.format",
              "type": "string",
              "value": "otlp"
            }
          ],
          "logs": [
            
          ],
          "processID": "p1",
          "warnings": null
        },
        {
          "traceID": "43337abeaf760fe1ef2a02704fc7c12c",
          "spanID": "18c3d5ae4a7c587b",
          "operationName": "被调2",
          "references": [
            {
              "refType": "CHILD_OF",
              "traceID": "43337abeaf760fe1ef2a02704fc7c12c",
              "spanID": "942af23b7b7a9ba1"
            }
          ],
          "startTime": 1720166470128076,
          "duration": 1000063,
          "tags": [
            {
              "key": "otel.library.name",
              "type": "string",
              "value": "主调微服务"
            },
            {
              "key": "span.kind",
              "type": "string",
              "value": "internal"
            },
            {
              "key": "internal.span.format",
              "type": "string",
              "value": "otlp"
            }
          ],
          "logs": [
            
          ],
          "processID": "p1",
          "warnings": null
        }
      ],
      "processes": {
        "p1": {
          "serviceName": "主调微服务名",
          "tags": [
            {
              "key": "deployment.environment",
              "type": "string",
              "value": "production"
            },
            {
              "key": "deployment.namespace",
              "type": "string",
              "value": "广州"
            },
            {
              "key": "service.instance.id",
              "type": "string",
              "value": "主调微服务名"
            },
            {
              "key": "service.instance.ip",
              "type": "string",
              "value": "1.2.3.4"
            }
          ]
        }
      },
      "warnings": null
    }
  ],
  "total": 0,
  "limit": 0,
  "offset": 0,
  "errors": null
}
```