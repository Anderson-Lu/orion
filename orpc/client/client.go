package client

import (
	"context"
	"errors"

	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/orpc/client/resolver"
	"github.com/Anderson-Lu/orion/orpc/codes"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
	"google.golang.org/grpc"
)

var (
	_OrionDefaultClientConfig = &OrionClientConfig{DailTimeout: 1000}
)

type OrionClientConfig struct {
	// server
	Host string

	// timout milliseconds for async dialing
	DailTimeout int64
}

func New(c *OrionClientConfig) (*OrionClient, error) {
	if c == nil {
		c = _OrionDefaultClientConfig
	}

	cli := &OrionClient{c: c}
	cli.breaker = circuit_break.NewCircuitBreaker()
	cli.oc = resolver.NewDefaultResolver(c.Host)

	return cli, nil
}

type OrionClient struct {
	c       *OrionClientConfig
	oc      resolver.IResolver
	breaker *circuit_break.CircuitBreaker
}

func (o *OrionClient) Invoke(ctx context.Context, method string, req, rsp interface{}, opts ...options.OrionClientInvokeOption) error {

	oCtx := newContext(ctx, opts...)
	if err := o.checkBreak(oCtx, method); err != nil {
		return o.after(oCtx, method, req, rsp)
	}

	conn, err := o.checkResolver(oCtx)
	if err != nil {
		return o.after(oCtx, method, req, rsp)
	}

	oCtx.err = conn.Invoke(oCtx.ctx, method, req, rsp, oCtx.options()...)
	return o.after(oCtx, method, req, rsp)
}

func (o *OrionClient) checkResolver(ctx *Context) (*grpc.ClientConn, error) {
	conn, err := o.oc.Select()
	if err != nil {
		ctx.err = err
		return nil, err
	}
	return conn, nil
}

func (o *OrionClient) checkBreak(ctx *Context, method string) error {

	if ok, ruleConfig := ctx.matchBreaker(method); o.breaker != nil && ok && ruleConfig != nil {
		o.breaker.Register(ruleConfig)
		if canPass := o.breaker.Pass(method); !canPass {
			return codes.WrapCodeFromError(errors.New("circuit break"), codes.ErrCodeCircuitBreak)
		}
	}
	return nil
}

func (o *OrionClient) after(ctx *Context, method string, req, rsp interface{}) error {
	reqCost := ctx.cost()
	if ok, _ := ctx.matchBreaker(method); ok && o.breaker != nil && codes.GetCodeFromError(ctx.err) != codes.ErrCodeCircuitBreak {
		o.breaker.Report(method, ctx.err == nil, int64(reqCost))
	}
	return ctx.err
}
