package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client"
	"github.com/Anderson-Lu/orion/orpc/client/balancer"
	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/orpc/client/resolver"
	"github.com/Anderson-Lu/orion/orpc/tracing"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
	"go.opentelemetry.io/otel/trace"
)

func main() {

	// cli, err := client.New(resolver.NewDirectResolver("127.0.0.1:8080"))

	crc := balancer.NewCrc32HashBalancer()
	rsv := resolver.NewConsulResovler("127.0.0.1:8500", resolver.WithBalancer(crc))

	cli, err := client.New(rsv)
	if err != nil {
		panic(err)
	}
	cli.RegisterCircuitBreakRule(&circuit_break.RuleConfig{
		Name:             "/todo.UitTodo/Add",
		Window:           &circuit_break.WindowConfig{Duration: 1000, Buckets: 10},
		OpenDuration:     10000,
		HalfOpenDuration: 10000,
		HaflOpenPassRate: 0,
		RuleExpression:   "req_count > 100",
	})

	bs := &tracing.Resources{}
	bs.SetEnv("test")
	bs.SetIP("9.134.188.178")
	bs.SetInstanceId("local")
	bs.SetNamespace("gz")
	bs.SetServiceName("10-30-client")

	tr, err := tracing.NewTracing(tracing.WithOpenTelemetryAddress("127.0.0.1:4317"), tracing.WithResource(bs))
	if err != nil {
		panic(err)
	}

	tr.Start()
	defer tr.Shutdown(context.Background())

	cli.RegisterTracing(tr)
	if err != nil {
		panic(err)
	}

	var rootSpan trace.Span
	rootCtx, rootSpan := tr.InternalSpan(context.Background(), "client-root-span")
	defer rootSpan.End()

	opts := []options.OrionClientInvokeOption{
		options.WithJson(),
		options.WithCircuitBreak(),
		options.WithService("mine.namespace.demo"),
		options.WithDirectAddress("127.0.0.1:8080"),
		options.WithMethod("/todo.UitTodo/Add"),
		options.WithBalancerParams("xxx"),
		options.WithHeaders("k1", "v1"),
	}

	req := &AddReq{Item: &TodoItem{Id: "1"}}
	rsp := &AddRsp{}
	orionRequest := client.NewOrionRequest(rootCtx, req, rsp, opts...)
	err = cli.Do(orionRequest)
	fmt.Println("rsp", rsp, "err", err)

	_, span2 := cli.GetTracing().InternalSpan(rootCtx, "client handle")
	time.Sleep(time.Millisecond * 100)
	span2.End()

	time.Sleep(time.Millisecond * 1000)

	for i := 0; i < 10; i++ {
		cli.GetTracing().Metrics().Counter(context.Background(), "count1", int64(i+1), "demo counter")
		cli.GetTracing().Metrics().Counter(context.Background(), "count2", int64(i+1), "demo counter")
	}

	select {}
}
