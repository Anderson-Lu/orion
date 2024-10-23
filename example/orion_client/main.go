package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client"
	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
)

func main() {

	cli, err := client.New(&client.OrionClientConfig{Host: "127.0.0.1:8080", DailTimeout: 10000})
	cli.RegisterCircuitBreakRule(&circuit_break.RuleConfig{
		Name:             "/todo.UitTodo/Add",
		Window:           &circuit_break.WindowConfig{Duration: 1000, Buckets: 10},
		OpenDuration:     10000,
		HalfOpenDuration: 10000,
		HaflOpenPassRate: 0,
		RuleExpression:   "req_count > 0",
	})
	if err != nil {
		panic(err)
	}

	opts := []options.OrionClientInvokeOption{
		options.WithJson(),
		options.WithCircuitBreak(),
	}

	for i := 0; i < 1000; i++ {
		req := &AddReq{Item: &TodoItem{Id: "1"}}
		rsp := &AddRsp{}
		err := cli.Invoke(context.Background(), "/todo.UitTodo/Add", req, rsp, opts...)
		time.Sleep(time.Millisecond * 300)
		fmt.Println("rsp", rsp, "err", err)
	}
}
