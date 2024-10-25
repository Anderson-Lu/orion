package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	ost "go.opentelemetry.io/otel/sdk/trace"
	ot "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

type Tracing struct {
	tracer    ot.Tracer
	provider  *ost.TracerProvider
	generator any
	processor any
	exportor  ost.SpanExporter
	err       error
	name      string
}

func NewTracing(name string) (*Tracing, error) {
	p := &Tracing{}
	if p.exportor, p.err = stdouttrace.New(stdouttrace.WithPrettyPrint()); p.err != nil {
		return nil, p.err
	}
	p.tracer = otel.Tracer(name)
	p.provider = ost.NewTracerProvider(ost.WithBatcher(p.exportor, ost.WithBatchTimeout(time.Second)))
	p.name = name
	return p, nil
}

func (p *Tracing) Start() {
	otel.SetTracerProvider(p.provider)
}

func (p *Tracing) Shutdown(ctx context.Context) {
	p.provider.Shutdown(ctx)
}

func (p *Tracing) Span(ctx context.Context, spanName string) (context.Context, ot.Span) {
	return p.tracer.Start(ctx, spanName, ot.WithAttributes(attribute.KeyValue{
		Key:   "service.name",
		Value: attribute.StringValue(p.name),
	}))
}
