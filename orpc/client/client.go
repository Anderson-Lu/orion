package client

import (
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

func (o *OrionClient) GetTracing() *tracing.Tracing {
	return o.trace
}

func (o *OrionClient) RegisterCircuitBreakRule(ruleConfigs ...*circuit_break.RuleConfig) {
	if o.breaker == nil {
		return
	}
	for _, v := range ruleConfigs {
		o.breaker.Register(v)
	}
}

func (o *OrionClient) Do(request *OrionRequest) error {

	if o.trace != nil {
		request.ctx, request.span = o.trace.OutgoingSpan(request.ctx, "orion-grpc:"+request.method, request.metaCarrier())
	}

	if circuitKey := request.getCircuitKey(); o.breaker != nil && circuitKey != "" {
		if canPass := o.breaker.Pass(circuitKey); !canPass {
			request.wrapError(codes.ErrClientCircuitBreaked)
			return o.after(request)
		}
	}

	var conn *grpc.ClientConn
	var err error
	if request.directEnable {
		drsv := resolver.NewDirectResolver(request.direct)
		conn, err = drsv.Select(request.resolverKey, request.balancerParams...)
	} else {
		conn, err = o.rsv.Select(request.resolverKey, request.balancerParams...)
	}

	if err != nil {
		request.wrapError(err)
		return o.after(request)
	}

	request.wrapError(conn.Invoke(request.buildContext(), request.method, request.req, request.rsp, request.callOptions...))
	return o.after(request)
}

func (o *OrionClient) after(request *OrionRequest) error {
	reqCost := request.cost()
	if circuitKey := request.getCircuitKey(); o.breaker != nil && circuitKey != "" && codes.GetCodeFromError(request.err()) != codes.ErrCodeCircuitBreak {
		o.breaker.Report(circuitKey, len(request.errs) == 0, int64(reqCost))
	}

	if request.span != nil {
		code, msg := codes.GetCodeAndMessageFromError(request.err())
		if code != 0 {
			request.span.SetStatus(oCodes.Error, msg)
		}
		request.span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key(tracing.KEY_SPAN_ERRCODE),
			Value: attribute.IntValue(code),
		}, attribute.KeyValue{
			Key:   attribute.Key(tracing.KEY_UNI_TRACE_ID),
			Value: attribute.StringValue(request.span.SpanContext().TraceID().String()),
		})
		request.span.End()
	}

	return request.err()
}
