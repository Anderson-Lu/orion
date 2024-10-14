package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client"
)

func main() {
	cli, err := client.New(&client.OrionClientConfig{
		Host:               "127.0.0.1:8080",
		DailTimeout:        3000,
		ConnectionNum:      1,
		ConnectionBalancer: "json",
	})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &AddReq{Item: &TodoItem{Id: "1"}}
	rsp := &AddRsp{}

	if err := cli.Invoke(ctx, "/todo.UitTodo/Add", req, rsp, client.WithJson(), client.WithHash("uid")); err != nil {
		panic(err)
	}

	fmt.Println("rsp", rsp, "err", err)
}
