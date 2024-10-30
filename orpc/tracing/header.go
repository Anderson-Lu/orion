package tracing

import (
	"context"

	ot "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	KEY_HEADER_TRACE_ID = "x-orion-traceid"
	KEY_HEADER_SPAN_ID  = "x-orion-spanid"
)

func NewTraceContext(ctx context.Context) *TraceContext {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &TraceContext{}
	}

	tc := &TraceContext{}
	if traceId := md.Get(KEY_HEADER_TRACE_ID); ok && len(traceId) > 0 {
		tc.TraceId = traceId[0]
	}
	if spanId := md.Get(KEY_HEADER_SPAN_ID); ok && len(spanId) > 0 {
		tc.SpanId = spanId[0]
	}
	return tc
}

type TraceContext struct {
	TraceId string
	SpanId  string
}

func (tc *TraceContext) ToSpanContext(ctx context.Context) context.Context {
	if len(tc.TraceId) < 16 || len(tc.SpanId) < 8 {
		return context.Background()
	}
	ox := ot.NewSpanContext(ot.SpanContextConfig{
		TraceID: ot.TraceID([]byte(tc.TraceId)),
		SpanID:  ot.SpanID([]byte(tc.SpanId)),
	})
	return ot.ContextWithSpanContext(ctx, ox)
}
