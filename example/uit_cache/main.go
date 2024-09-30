package main

import (
	"context"
	"fmt"
	"time"

	"github.com/uit/pkg/utils/ticker_cache"
)

func main() {

	tc := ticker_cache.New(
		ticker_cache.WithOrigin(&BarCache{}),
		ticker_cache.WithTTL(10),
		ticker_cache.WithElimination(60),
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
