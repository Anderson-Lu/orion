package client

import (
	"context"
	"fmt"

	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/orpc/client/resolver"
	"github.com/Anderson-Lu/orion/orpc/codes"
	"github.com/Anderson-Lu/orion/orpc/tracing"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
	"go.opentelemetry.io/otel/attribute"
	oCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
)

func New(rsv resolver.IResolver) (*OrionClient, error) {
	cli := &OrionClient{}
	cli.breaker = circuit_break.NewCircuitBreaker()
	cli.rsv = rsv

	return cli, nil
}

type OrionClient struct {
	rsv     resolver.IResolver
	breaker *circuit_break.CircuitBreaker
	trace   *tracing.Tracing
}

func (o *OrionClient) RegisterTracing(trace *tracing.Tracing) {
	o.trace = trace
}

func (o *OrionClient) RegisterCircuitBreakRule(ruleConfigs ...*circuit_break.RuleConfig) {
	if o.breaker == nil {
		return
	}
	for _, v := range ruleConfigs {
		o.breaker.Register(v)
	}
}

func (o *OrionClient) Invoke(ctx context.Context, req, rsp interface{}, opts ...options.OrionClientInvokeOption) error {

	meta := newOrionRequestMeta(ctx, req, rsp, opts...)
	if o.trace != nil {
		meta.ctx, meta.span = o.trace.SpanClient(ctx, meta.method)
		meta.headers.Set(tracing.KEY_HEADER_TRACE_ID, meta.span.SpanContext().TraceID().String())
	}

	if circuitKey := meta.getCircuitKey(); o.breaker != nil && circuitKey != "" {
		if canPass := o.breaker.Pass(circuitKey); !canPass {
			meta.wrapError(codes.ErrClientCircuitBreaked)
			return o.after(meta)
		}
	}

	var conn *grpc.ClientConn
	var err error
	if meta.directEnable {
		drsv := resolver.NewDirectResolver(meta.direct)
		conn, err = drsv.Select(meta.resolverKey, meta.balancerParams...)
	} else {
		conn, err = o.rsv.Select(meta.resolverKey, meta.balancerParams...)
	}

	if err != nil {
		meta.wrapError(err)
		return o.after(meta)
	}

	meta.wrapError(conn.Invoke(meta.buildContext(), meta.method, req, rsp, meta.callOptions...))
	return o.after(meta)
}

func (o *OrionClient) after(meta *OrionRequestMeta) error {
	reqCost := meta.cost()
	if circuitKey := meta.getCircuitKey(); o.breaker != nil && circuitKey != "" && codes.GetCodeFromError(meta.err()) != codes.ErrCodeCircuitBreak {
		o.breaker.Report(circuitKey, len(meta.errs) == 0, int64(reqCost))
	}
	code, msg := codes.GetCodeAndMessageFromError(meta.err())
	fmt.Println("code", code, "msg", msg)
	if code != 0 {
		meta.span.SetStatus(oCodes.Error, msg)
		meta.span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key(tracing.KEY_SPAN_ERRCODE),
			Value: attribute.IntValue(code),
		})
	}
	meta.span.End()
	return meta.err()
}
