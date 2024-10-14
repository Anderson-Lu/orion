package client

import (
	"context"
	"errors"

	"github.com/Anderson-Lu/orion/pkg/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	_OrionDefaultClientConfig = &OrionClientConfig{
		DailTimeout:   1000,
		ConnectionNum: 1,
	}
)

type OrionClientConfig struct {
	// server
	Host string

	// timout milliseconds for async dialing
	DailTimeout int64

	// if set to zero, ConnectionNum = 1 will be set.
	ConnectionNum int64

	// can be set to "random" or "hash", if not set, random balancer will be set by default.
	// it only works when ConnectionNum > 0
	ConnectionBalancer string
}

func New(c *OrionClientConfig) (*OrionClient, error) {
	if c == nil {
		c = _OrionDefaultClientConfig
	}
	if c.ConnectionNum <= 0 {
		c.ConnectionNum = 1
	}
	if c.ConnectionBalancer == "" {
		c.ConnectionBalancer = "random"
	}
	cli := &OrionClient{c: c}
	cli.initBalancer()
	for i := 0; i < int(c.ConnectionNum); i++ {
		c, err := grpc.NewClient(c.Host, cli.dailOptions()...)
		if err != nil {
			return nil, err
		}
		cli.conns = append(cli.conns, c)
	}

	return cli, nil
}

type OrionClient struct {
	b     balancer.Balancer
	c     *OrionClientConfig
	conns []grpc.ClientConnInterface
}

func (o *OrionClient) initBalancer() {
	switch o.c.ConnectionBalancer {
	case "random":
		o.b = balancer.NewRandomBalancer(int(o.c.ConnectionNum))
	case "hash":
		o.b = balancer.NewHashBalancer(int(o.c.ConnectionNum))
	default:
		o.b = balancer.NewDefaultBalancer(int(o.c.ConnectionNum))
	}
}

func (o *OrionClient) dailOptions() []grpc.DialOption {
	r := []grpc.DialOption{}
	r = append(r, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return r
}

func (o *OrionClient) Invoke(ctx context.Context, method string, req, rsp interface{}, opts ...OrionClientInvokeOption) error {
	if len(o.conns) == 0 {
		return errors.New("nil conn")
	}
	bIdx := o.b.Get(o.balanceOptions(opts...))
	defer o.b.Update(bIdx)
	options := o.buildOptions(opts...)
	return o.conns[bIdx].Invoke(ctx, method, req, rsp, options...)
}

func (o *OrionClient) balanceOptions(opts ...OrionClientInvokeOption) string {
	for _, opt := range opts {
		switch opt.Type() {
		case OptionTypeBalanceOption:
			for _, v := range opt.Params() {
				return v.(string)
			}
		}
	}
	return ""
}

func (o *OrionClient) buildOptions(opts ...OrionClientInvokeOption) []grpc.CallOption {
	r := []grpc.CallOption{}
	for _, opt := range opts {
		switch opt.Type() {
		case OptionTypeGrpcCallOption:
			for _, v := range opt.Params() {
				r = append(r, v.(grpc.CallOption))
			}
		}
	}
	return r
}
