package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
)

func main() {

	cli, err := client.New(&client.OrionClientConfig{
		Host:               "127.0.0.1:8080",
		DailTimeout:        3000,
		ConnectionNum:      1,
		ConnectionBalancer: "json",
		CircuitBreakRules: []*circuit_break.RuleConfig{
			{
				Name:             "/todo.UitTodo/Add",
				Window:           &circuit_break.WindowConfig{},
				OpenDuration:     1000,
				HalfOpenDuration: 100,
				HaflOpenPassRate: 0,
				RuleExpression:   "req_count >= 1 && succ_rate < 0.90",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	req := &AddReq{Item: &TodoItem{Id: "1"}}
	rsp := &AddRsp{}

	opts := []client.OrionClientInvokeOption{
		client.WithJson(),
		client.WithHash("uid"),
		client.WithCircuitBreak("/todo.UitTodo/Add"),
	}

	for i := 0; i < 1000; i++ {
		err := cli.Invoke(context.Background(), "/todo.UitTodo/Add", req, rsp, opts...)
		time.Sleep(time.Millisecond * 300)
		fmt.Println("rsp", rsp, "err", err)
	}
}
