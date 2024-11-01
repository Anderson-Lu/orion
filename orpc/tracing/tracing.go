package tracing

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/sdk/trace"
	ot "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

type Tracing struct {
	name   string
	tracer ot.Tracer

	exportorAddr string
	exportorConn *grpc.ClientConn

	metricExportor *otlpmetricgrpc.Exporter
	metricProvider *metric.MeterProvider
	metrics        *OrionMetrics

	traceExportor *otlptrace.Exporter
	traceProvider *trace.TracerProvider

	baseResources *Resources

	propagator OrionTracingPropagation
}

func NewTracing(opts ...TracingOption) (*Tracing, error) {
	p := &Tracing{}

	for _, v := range opts {
		v(p)
	}
	if err := p.initExportor(); err != nil {
		return nil, err
	}
	if err := p.initProvider(); err != nil {
		return nil, err
	}
	p.initTracer()
	p.initMetrics()
	return p, nil
}

func (p *Tracing) initExportor() error {
	if p.exportorAddr == "" {
		return errors.New("empty exportorAddr specified")
	}

	conn, err := grpc.NewClient(p.exportorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	p.exportorConn = conn

	te, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithGRPCConn(p.exportorConn))
	if err != nil {
		return err
	}
	p.traceExportor = te

	me, err := otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithGRPCConn(p.exportorConn))
	if err != nil {
		return err
	}
	p.metricExportor = me

	return nil
}

func (p *Tracing) initProvider() error {
	rso, err := resource.New(context.Background(), resource.WithAttributes(p.baseResources.kvs...))
	if err != nil {
		return err
	}
	traceProcessor := trace.NewBatchSpanProcessor(p.traceExportor)
	p.traceProvider = trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()), trace.WithSpanProcessor(traceProcessor), trace.WithResource(rso))
	p.metricProvider = metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(p.metricExportor)), metric.WithResource(rso))
	otel.SetTracerProvider(p.traceProvider)
	otel.SetMeterProvider(p.metricProvider)
	return nil
}

func (p *Tracing) initMetrics() {
	p.metrics = NewOrionMetrics(p.name, p.metricProvider)
}

func (p *Tracing) initTracer() error {
	p.propagator = OrionTracingPropagation{}
	p.tracer = p.traceProvider.Tracer(p.name)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(p.propagator))
	return nil
}

func (p *Tracing) Start() {}

func (p *Tracing) Shutdown(ctx context.Context) {
	p.traceProvider.Shutdown(ctx)
	p.metricProvider.Shutdown(ctx)
}

func (p *Tracing) OutgoingSpan(ctx context.Context, spanName string, mds MetadataCarrier, attrKvs ...string) (context.Context, ot.Span) {
	tracerName := mds.SetTracerNameNx(spanName)
	ctx = p.propagator.WrapContextWithTracerName(ctx, tracerName)
	attrKvs = append(attrKvs, KEY_RESOURCE_SERVICE_NAME, tracerName)
	ctx, span := p.Span(ctx, spanName, ot.SpanKindClient, attrKvs...)
	otel.GetTextMapPropagator().Inject(ctx, propagation.TextMapCarrier(mds))
	return ctx, span
}

func (p *Tracing) ServerSpan(ctx context.Context, spanName string, attrKvs ...string) (context.Context, ot.Span) {
	ctx, span := p.Span(ctx, spanName, ot.SpanKindServer, attrKvs...)
	return ctx, span
}

func (p *Tracing) InternalSpan(ctx context.Context, spanName string, attrKvs ...string) (context.Context, ot.Span) {
	ctx, span := p.Span(ctx, spanName, ot.SpanKindInternal, attrKvs...)
	return ctx, span
}

func (p *Tracing) Span(ctx context.Context, spanName string, spanKind ot.SpanKind, attrKvs ...string) (context.Context, ot.Span) {

	kvs := []attribute.KeyValue{}
	if len(attrKvs)%2 == 0 {
		for i := 0; i < len(attrKvs); i += 2 {
			kvs = append(kvs, attribute.KeyValue{
				Key:   attribute.Key(attrKvs[0]),
				Value: attribute.StringValue(attrKvs[1]),
			})
		}
	}

	return otel.Tracer(p.propagator.GetTracerName(ctx)).Start(ctx, spanName, ot.WithAttributes(kvs...), ot.WithSpanKind(spanKind))
}

func (p *Tracing) Metrics() *OrionMetrics {
	return p.metrics
}
