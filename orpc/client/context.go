package client

import (
	"context"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client/options"
	"google.golang.org/grpc"
)

type Context struct {
	ctx   context.Context
	begin int64
	end   int64
	err   error

	opts []options.OrionClientInvokeOption
}

func (c *Context) matchBreaker() bool {
	for _, v := range c.opts {
		if v.Type() == options.OptionTypeCircuitBreakOption {
			return true
		}
	}
	return false
}

func (c *Context) cost() int64 {
	c.end = time.Now().UnixMilli()
	return c.end - c.begin
}

func (c *Context) options() []grpc.CallOption {
	r := []grpc.CallOption{}
	for _, v := range c.opts {
		switch k := v.(type) {
		case *options.CallOptionWithJson:
			r = append(r, k.GrpcCallOption())
		}
	}
	return r
}

func newContext(ctx context.Context, opts ...options.OrionClientInvokeOption) *Context {
	return &Context{
		ctx:   ctx,
		begin: time.Now().UnixMilli(),
		opts:  opts,
	}
}
