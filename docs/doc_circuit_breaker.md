# 术语定义

**Rule(规则)**

OrionCircuitBreaker支持定义多组规则, 通过规则的资源名字(ResourceName)作为唯一标识.

**Window(窗口)**

要判断是否触发熔断,需要在一定的时间段内收集指标数据,这个时间段称之为**Window**.如以10s作为一个滑动窗口,规则判断则基于当前时间点所在的窗口来判断.

**Bucket(桶)**

一个窗口下分为多个桶,用来存储窗口下的每个时刻的指标数据,如100ms一个桶.

**Stratege(策略)**

当触发熔断后,需要进行一定的拦截操作,比如全量请求拦截、部分请求放行等,此时需要有对应的策略管理器来执行熔断器的状态变更以及柔性降级等操作.

# 配置项

```go
type RuleConfig struct {

// the name of the rule.
  Name string
  
  // specify the stat window including buckets size and duration.
  Window *WindowConfig
  
  OpenDuration     int64 // if circuit-break occours, it will recover after next BreakDuration duration
  HalfOpenDuration int64 // eg: 60s = 60 * 1000 = 60000
  HaflOpenPassRate int64 // 0-100

  // specify different rule checker to handle circuit break.
  BreakHandlerType string

  // expression for checking,
  // suuport ops: [ +, -, *, /, &&, || ]
  // support vars: [req_count, req_succ_count, req_fail_count, succ_rate]
  // eg: `req_count > 100 && succ_rate < 0.95`
  RuleExpression string
}
```

# 熔断表达式

OrionCircuitBreaker支持业务通过表达式实现灵活的熔断策略, 目前支持的字段如下:

```go
  FieldRequestNum     string = "req_count"
  FieldRequestSuccNum string = "req_succ_count"
  FieldRequestFailNum string = "req_fail_count"
  FieldSuccRate       string = "succ_rate"
  FieldAvgCost        string = "avg_cost"
```

举例: 在**当前窗口内请求数>100且成功率低于50%**时触发熔断, 此时可以采用表达式:

```go
req_count > 100 && succ_rate < 0.50
```

Orion内部会通过解析AST来判断表达式是否成立,从而决定熔断器是否发生状态流转.

# 熔断器状态流转

```go
const (
  CircuitBreakStatusClose    CircuitBreakStatus = 0
  CircuitBreakStatusHalfOpen CircuitBreakStatus = 1
  CircuitBreakStatusOpen     CircuitBreakStatus = 2
)
```

- `CircuitBreakStatusClose`     关闭状态, 当前状态下全部允许通过.
- `CircuitBreakStatusHalfOpen`  半开启状态, 当前状态下会有一定几率允许通过.
- `CircuitBreakStatusOpen`      开启状态, 当前状态下不允许通过.

1. 初始熔断器为**Close**状态, 当窗口规则判断异常时, **Close --> HalfOpen**状态, 否则依旧维持**Close**状态
2. 当状态为**HalfOpen**时, 如果在**HalfDuration**后检测到恢复, 则转为**Close**状态, 否则转向**Open**状态
3. 当状态为**Open**时, 如果在**OpenDuration**后检测恢复, 则转为**HalfOpen**状态, 否则依旧维持**Open**状态

# 快速开始

更多详情,参见:[容器器示例](../example/circuit_break/main.go)

```go
rc := &circuit_break.RuleConfig{
  Name: "rule1",
  Window: &circuit_break.WindowConfig{
      Duration: 10000,
      Buckets:  10,
    },
    OpenDuration:     1000,
    HalfOpenDuration: 1000,
    HaflOpenPassRate: 0,
    RuleExpression:   "req_count >= 1 && succ_rate < 0.90",
  }

cb := circuit_break.NewCircuitBreaker()
cb.Register(rc)
cb.Report(rc.Name, true, cost) // 上报结果
pass := cb.Pass(rc.Name) // 查询是否需要熔断
if pass {
    // 熔断了
}
```