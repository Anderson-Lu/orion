package main

import (
	"context"
	"time"

	"github.com/Anderson-Lu/orion/orpc/tracing"
	"go.opentelemetry.io/otel/codes"
)

func main() {

	ctx := context.Background()

	bs := &tracing.Resources{}
	bs.Env("production")
	bs.IP("9.134.188.178")
	bs.InstanceId("主调微服务名")
	bs.ServiceName("主调微服务名")
	bs.Namespace("广州")

	p, err := tracing.NewTracing("主调微服务",
		tracing.WithOpenTelemetryAddress("127.0.0.1:4317"),
		tracing.WithResource(bs),
	)
	if err != nil {
		panic(err)
	}
	p.Start()
	defer p.Shutdown(ctx)

	ctx1, span1 := p.SpanClient(ctx, "precheck", "callee.service", "被调1")
	time.Sleep(time.Second)
	span1.SetStatus(codes.Error, "error occuro, 遇到错误了")
	span1.End()

	_, span2 := p.SpanClient(ctx1, "被调2")
	time.Sleep(time.Second)
	span2.End()

	time.Sleep(time.Second * 190)

}
