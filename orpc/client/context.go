package client

import (
	"context"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client/options"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
	"google.golang.org/grpc"
)

type Context struct {
	ctx   context.Context
	begin int64
	end   int64
	err   error

	opts []options.OrionClientInvokeOption
}

func (c *Context) matchBreaker(method string) (bool, *circuit_break.RuleConfig) {
	for _, v := range c.opts {
		if v.Type() == options.OptionTypeCircuitBreakOption {
			if matched, ruleConfig := (v.(*options.CircuitBreakerOption)).IsMatch(method); matched {
				return true, ruleConfig
			}
		}
	}
	return false, nil
}

func (c *Context) cost() int64 {
	c.end = time.Now().Unix()
	return c.end - c.begin
}

func (c *Context) options() []grpc.CallOption {
	r := []grpc.CallOption{}
	return r
}

func newContext(ctx context.Context, opts ...options.OrionClientInvokeOption) *Context {
	return &Context{
		ctx:   ctx,
		begin: time.Now().UnixMilli(),
		opts:  opts,
	}
}
