# 本地缓存(TickerCache)

一些本地需要初步缓存的场景,比如某些业务配置(强实时性不高),可以先在本地缓存一段时间后再刷新.在这种模式下,微服务场景可能涉及多个业务,多个组件都有相似的需求,因此,UIT统一封装了ticker_cache工具类,精简业务实现,避免各个业务自立山头.

# 快速开始

业务侧需要实现源数据的获取逻辑,比如从远程`redis`/`mysql`等存储介质中拉取数据,只需要实现以下接口:

```go
type IFetcher interface {
  FetchFromOrigin(ctx context.Context, key string) (interface{}, error)
}
```

然后将其作为参数传递给`ticker_cache`组件即可实现本地缓存的功能:

```go
package main

import (
  "context"
  "fmt"
  "time"

  "github.com/orion/pkg/utils/ticker_cache"
)

func main() {

  tc := ticker_cache.New(
    ticker_cache.WithOrigin(&BarCache{}), // 数据源刷新实现
    ticker_cache.WithTTL(10),             // 缓存时长
    ticker_cache.WithElimination(60),     // 惰性缓存清理(可选)
  )

  fooValue, err := tc.Get(context.Background(), "foo")
  fmt.Println("fooValue", fooValue, "err", err)

  barValue, err := tc.Get(context.Background(), "bar")
  fmt.Println("barValue", barValue, "err", err)

  time.Sleep(time.Second * 3)

  fooValue, err = tc.Get(context.Background(), "foo")
  fmt.Println("fooValue", fooValue, "err", err)

  barValue, err = tc.Get(context.Background(), "bar")
  fmt.Println("barValue", barValue, "err", err)

}

type BarCache struct{}

func (b *BarCache) FetchFromOrigin(ctx context.Context, key string) (interface{}, error) {
  // use your logic, for example, fetch data from redis, mysql or any other middlewares.
  return fmt.Sprintf("key [%s] at [%d]", key, time.Now().Unix()), nil
}
```