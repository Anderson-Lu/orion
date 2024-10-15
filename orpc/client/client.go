package client

import (
	"context"
	"errors"
	"time"

	"github.com/Anderson-Lu/orion/orpc/codes"
	"github.com/Anderson-Lu/orion/pkg/balancer"
	"github.com/Anderson-Lu/orion/pkg/circuit_break"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	_OrionDefaultClientConfig = &OrionClientConfig{
		DailTimeout:   1000,
		ConnectionNum: 1,
	}
	_OrionReqBeginTimeKey = struct{}{}
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

	// if set, circuit breaker will be set
	CircuitBreakRules []*circuit_break.RuleConfig
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
	cli.initCircuitBreaker()

	return cli, nil
}

type OrionClient struct {
	b       balancer.Balancer
	c       *OrionClientConfig
	conns   []grpc.ClientConnInterface
	breaker *circuit_break.CircuitBreaker
}

func (o *OrionClient) initCircuitBreaker() {
	if len(o.c.CircuitBreakRules) <= 0 {
		return
	}
	o.breaker = circuit_break.NewCircuitBreaker()
	for _, rule := range o.c.CircuitBreakRules {
		o.breaker.Register(rule)
	}
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

	ctx, err := o.beforeInvoke(ctx, method, req, rsp, opts...)
	if err != nil {
		return o.afterInvoke(ctx, method, req, rsp, err, opts...)
	}

	options := o.buildOptions(opts...)
	err = o.conns[bIdx].Invoke(ctx, method, req, rsp, options...)
	return o.afterInvoke(ctx, method, req, rsp, err, opts...)
}

func (o *OrionClient) beforeInvoke(ctx context.Context, method string, req, rsp interface{}, opts ...OrionClientInvokeOption) (context.Context, error) {

	ctx = context.WithValue(ctx, _OrionReqBeginTimeKey, time.Now().UnixMilli())

	needCheckCircuit := false
	for _, v := range opts {
		if v.Type() == OptionTypeCircuitBreakOption {
			needCheckCircuit = true
		}
	}

	if needCheckCircuit && o.breaker != nil {
		if canPass := o.breaker.Pass(method); !canPass {
			return ctx, codes.WrapCodeFromError(errors.New("circuit break"), codes.ErrCodeCircuitBreak)
		}
	}

	return ctx, nil
}

func (o *OrionClient) afterInvoke(ctx context.Context, method string, req, rsp interface{}, err error, opts ...OrionClientInvokeOption) error {

	needCheckCircuit := false
	for _, v := range opts {
		if v.Type() == OptionTypeCircuitBreakOption {
			needCheckCircuit = true
		}
	}

	reqCost := o.cost(ctx)

	if needCheckCircuit && o.breaker != nil && codes.GetCodeFromError(err) != codes.ErrCodeCircuitBreak {
		o.breaker.Report(method, err == nil, int64(reqCost))
	}

	return err
}

func (o *OrionClient) cost(ctx context.Context) int64 {
	begin := ctx.Value(_OrionReqBeginTimeKey)
	if begin, ok := begin.(int64); ok {
		return (time.Now().UnixMilli() - begin)
	}
	return 0.0
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
