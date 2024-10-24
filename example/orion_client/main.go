package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client"
	"github.com/Anderson-Lu/orion/orpc/client/balancer"
	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/orpc/client/resolver"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
)

func main() {

	// cli, err := client.New(resolver.NewDirectResolver("127.0.0.1:8080"))

	crc := balancer.NewCrc32HashBalancer()
	rsv := resolver.NewConsulResovler("127.0.0.1:8500", resolver.WithBalancer(crc))
	cli, err := client.New(rsv)
	cli.RegisterCircuitBreakRule(&circuit_break.RuleConfig{
		Name:             "/todo.UitTodo/Add",
		Window:           &circuit_break.WindowConfig{Duration: 1000, Buckets: 10},
		OpenDuration:     10000,
		HalfOpenDuration: 10000,
		HaflOpenPassRate: 0,
		RuleExpression:   "req_count > 100",
	})
	if err != nil {
		panic(err)
	}

	opts := []options.OrionClientInvokeOption{
		options.WithJson(),
		options.WithCircuitBreak(),
		options.WithService("mine.namespace.demo"),
		options.WithDirectAddress("127.0.0.1:8080"),
		options.WithMethod("/todo.UitTodo/Add"),
		options.WithBalancerParams("xxx"),
		options.WithHeaders("k1", "v1"),
	}

	for i := 0; i < 1000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		req := &AddReq{Item: &TodoItem{Id: "1"}}
		rsp := &AddRsp{}
		err := cli.Invoke(ctx, req, rsp, opts...)
		cancel()
		time.Sleep(time.Millisecond * 300)
		fmt.Println("rsp", rsp, "err", err)
	}
}
