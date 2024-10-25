package main

import (
	"context"
	"time"

	"github.com/Anderson-Lu/orion/orpc/tracing"
)

func main() {

	ctx := context.Background()

	p, _ := tracing.NewTracing("mine.example.tracing")
	p.Start()
	defer p.Shutdown(ctx)

	ctx1, span1 := p.Span(ctx, "span1")
	time.Sleep(time.Second)
	span1.End()

	_, span2 := p.Span(ctx1, "span2")
	time.Sleep(time.Second)
	span2.End()

}
