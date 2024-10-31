package tracing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/multierr"
)

type OrionTracingPropagation struct{}

var _ propagation.TextMapPropagator = OrionTracingPropagation{}
var errInvalidSampledHeader = errors.New("invalid Sampled header found")
var errInvalidTraceIDHeader = errors.New("invalid TraceID header found")
var errInvalidSpanIDHeader = errors.New("invalid SpanID header found")
var errInvalidScope = errors.New("require either both traceID and spanID or none")
var emptySpanContext = trace.SpanContext{}

const (
	headerKeyTraceId       = "x-orion-tracer-traceid"
	headerKeySpanId        = "x-orion-tracer-spanid"
	headerKeySampled       = "x-orion-tracer-sampled"
	headerKeyTracerName    = "x-orion-tracer-name"
	headerKeyBaggagePrefix = "x-orion-tracer-baggage-"

	traceIdPadding     = "0000000000000000"
	traceID64BitsWidth = 64 / 2 // 16 hex character Trace ID.
)

func (o OrionTracingPropagation) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
		return
	}
	carrier.Set(headerKeyTraceId, sc.TraceID().String()[len(sc.TraceID().String())-traceID64BitsWidth:])
	carrier.Set(headerKeySpanId, sc.SpanID().String())
	if sc.IsSampled() {
		carrier.Set(headerKeySampled, "true")
	} else {
		carrier.Set(headerKeySampled, "false")
	}
	for _, m := range baggage.FromContext(ctx).Members() {
		carrier.Set(fmt.Sprintf("%s%s", headerKeyBaggagePrefix, m.Key()), m.Value())
	}
	carrier.Set(headerKeyTracerName, o.GetTracerName(ctx))
}

func (o OrionTracingPropagation) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {

	traceId := carrier.Get(headerKeyTraceId)
	spanId := carrier.Get(headerKeySpanId)
	sampled := carrier.Get(headerKeySampled)
	sc, err := o.extract(traceId, spanId, sampled)

	if err != nil || !sc.IsValid() {
		return ctx
	}

	if bags, err := o.extractBags(carrier); err == nil {
		ctx = baggage.ContextWithBaggage(ctx, bags)
	}

	tracerName := carrier.Get(headerKeyTracerName)
	ctx = o.WrapContextWithTracerName(ctx, tracerName)

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (o OrionTracingPropagation) Fields() []string {
	return []string{headerKeyTraceId, headerKeySpanId, headerKeySampled}
}

func (o OrionTracingPropagation) extract(traceID, spanID, sampled string) (trace.SpanContext, error) {
	var (
		err           error
		requiredCount int
		scc           = trace.SpanContextConfig{}
	)

	switch strings.ToLower(sampled) {
	case "0", "false", "1", "true", "":
	default:
		return emptySpanContext, errInvalidSampledHeader
	}

	if traceID != "" {
		requiredCount++
		id := traceID
		if len(traceID) == 16 {
			id = traceIdPadding + traceID
		}
		if scc.TraceID, err = trace.TraceIDFromHex(id); err != nil {
			return emptySpanContext, errInvalidTraceIDHeader
		}
	}

	if spanID != "" {
		requiredCount++
		if scc.SpanID, err = trace.SpanIDFromHex(spanID); err != nil {
			return emptySpanContext, errInvalidSpanIDHeader
		}
	}

	if requiredCount != 0 && requiredCount != 2 {
		return emptySpanContext, errInvalidScope
	}

	return trace.NewSpanContext(scc), nil
}

func (o OrionTracingPropagation) extractBags(carrier propagation.TextMapCarrier) (baggage.Baggage, error) {
	var err error
	var members []baggage.Member
	for _, key := range carrier.Keys() {
		lowerKey := strings.ToLower(key)
		if !strings.HasPrefix(lowerKey, headerKeyBaggagePrefix) {
			continue
		}
		strippedKey := strings.TrimPrefix(lowerKey, headerKeyBaggagePrefix)
		member, e := baggage.NewMember(strippedKey, carrier.Get(key))
		if e != nil {
			err = multierr.Append(err, e)
			continue
		}
		members = append(members, member)
	}
	bags, e := baggage.New(members...)
	if err != nil {
		return bags, multierr.Append(err, e)
	}
	return bags, err
}

type ContextTraceNameKey struct{}

func (o OrionTracingPropagation) WrapContextWithTracerName(ctx context.Context, traceName string) context.Context {
	if traceName == "" {
		return ctx
	}
	ctx = context.WithValue(ctx, ContextTraceNameKey{}, traceName)
	return ctx
}

func (o OrionTracingPropagation) GetTracerName(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(ContextTraceNameKey{})
	if v == nil {
		return ""
	}
	if c, ok := v.(string); ok {
		return c
	}
	return ""
}
