package client

import (
	"context"
	"time"

	"github.com/Anderson-Lu/orion/orpc/client/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type OrionRequestMeta struct {

	// origonal params
	ctx  context.Context
	opts []options.OrionClientInvokeOption

	// statistics
	begin int64
	end   int64
	errs  []error

	// discovery
	service        string
	method         string
	direct         string
	directEnable   bool
	resolverKey    string
	circuitEnable  bool
	circuitKey     string
	balancerParams []interface{}
	headers        metadata.MD

	// request
	req         interface{}
	rsp         interface{}
	callOptions []grpc.CallOption
}

func (c *OrionRequestMeta) err() error {
	if len(c.errs) == 0 {
		return nil
	}
	return c.errs[len(c.errs)-1]
}

func (c *OrionRequestMeta) cost() int64 {
	c.end = time.Now().UnixMilli()
	return c.end - c.begin
}

func (c *OrionRequestMeta) getCircuitKey() string {
	if !c.circuitEnable || c.circuitKey == "" {
		return ""
	}
	return c.circuitKey
}

func (c *OrionRequestMeta) wrapError(err error) {
	if err == nil {
		return
	}
	c.errs = append(c.errs, err)
}

func (c *OrionRequestMeta) buildContext() context.Context {
	if c.headers == nil {
		return c.ctx
	}

	cc := metadata.NewOutgoingContext(c.ctx, c.headers)
	return cc
}

func newOrionRequestMeta(ctx context.Context, req, rsp interface{}, opts ...options.OrionClientInvokeOption) *OrionRequestMeta {
	o := &OrionRequestMeta{
		ctx:   ctx,
		begin: time.Now().UnixMilli(),
		opts:  opts,
		req:   req,
		rsp:   rsp,
	}
	for _, v := range opts {
		switch opt := v.(type) {
		case *options.CallOptionWithService:
			o.service = opt.Service()
			o.direct = opt.Direct()
			if opt.IsDirect() {
				o.directEnable = true
			}
		case *options.CallOptionWithPath:
			o.method = opt.Method()
		case *options.CallOptionWithCircuitBreaker:
			o.circuitEnable = true
			o.circuitKey = opt.Key()
		case *options.CallOptionWithJson:
			o.callOptions = append(o.callOptions, opt.GrpcCallOption())
		case *options.CallOptionWithBalancerParams:
			o.balancerParams = v.Params()
		case *options.CallOptionWithHeader:
			o.headers = opt.Metadata()
		}
	}
	if o.service != "" {
		o.resolverKey = o.service
	} else {
		o.resolverKey = o.direct
	}
	return o
}
