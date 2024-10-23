package client

import (
	"context"
	"errors"

	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/orpc/client/resolver"
	"github.com/Anderson-Lu/orion/orpc/codes"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
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
}

func (o *OrionClient) RegisterCircuitBreakRule(ruleConfigs ...*circuit_break.RuleConfig) {
	if o.breaker == nil {
		return
	}
	for _, v := range ruleConfigs {
		o.breaker.Register(v)
	}
}

func (o *OrionClient) Invoke(ctx context.Context, method string, req, rsp interface{}, opts ...options.OrionClientInvokeOption) error {

	oCtx := newContext(ctx, opts...)
	if err := o.checkBreak(oCtx, method); err != nil {
		return o.after(oCtx, method, req, rsp)
	}

	conn, err := o.rsv.Select(method)
	if err != nil {
		return o.after(oCtx, method, req, rsp)
	}

	oCtx.err = conn.Invoke(oCtx.ctx, method, req, rsp, oCtx.options()...)
	return o.after(oCtx, method, req, rsp)
}

func (o *OrionClient) checkBreak(ctx *Context, method string) error {

	if ok := ctx.matchBreaker(); o.breaker != nil && ok {
		if canPass := o.breaker.Pass(method); !canPass {
			ctx.err = codes.WrapCodeFromError(errors.New("circuit break"), codes.ErrCodeCircuitBreak)
			return ctx.err
		}
	}
	return nil
}

func (o *OrionClient) after(ctx *Context, method string, req, rsp interface{}) error {
	reqCost := ctx.cost()
	if ok := ctx.matchBreaker(); ok && o.breaker != nil && codes.GetCodeFromError(ctx.err) != codes.ErrCodeCircuitBreak {
		o.breaker.Report(method, ctx.err == nil, int64(reqCost))
	}
	return ctx.err
}
