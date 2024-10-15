package client

import (
	"github.com/Anderson-Lu/orion/orpc/codec"
	"google.golang.org/grpc"
)

type OptionType int

const (
	OptionTypeGrpcCallOption     OptionType = 1
	OptionTypeBalanceOption      OptionType = 2
	OptionTypeCircuitBreakOption OptionType = 3
)

type OrionClientInvokeOption interface {
	Params() []interface{}
	Type() OptionType
}

type HashOption struct {
	p string
}

func WithHash(args string) OrionClientInvokeOption {
	return &HashOption{
		p: args,
	}
}

func (c HashOption) Type() OptionType {
	return OptionTypeBalanceOption
}

func (h HashOption) Params() []interface{} {
	return []interface{}{h.p}
}

func WithJson() OrionClientInvokeOption {
	return &CallOptionWithJson{}
}

type CallOptionWithJson struct {
}

func (c CallOptionWithJson) Params() []interface{} {
	return []interface{}{
		grpc.CallContentSubtype(codec.JSON{}.Name()),
	}
}

func (c CallOptionWithJson) Type() OptionType {
	return OptionTypeGrpcCallOption
}

func WithCircuitBreak(methods ...string) OrionClientInvokeOption {
	return &CircuitBreakerOption{methods: methods}
}

type CircuitBreakerOption struct {
	methods []string
}

func (c CircuitBreakerOption) Params() []interface{} {
	return []interface{}{c.methods}
}

func (c CircuitBreakerOption) Type() OptionType {
	return OptionTypeCircuitBreakOption
}
