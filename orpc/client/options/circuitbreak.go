package options

func WithCircuitBreak() OrionClientInvokeOption {
	return &CircuitBreakerOption{}
}

type CircuitBreakerOption struct {
}

func (c CircuitBreakerOption) Params() []interface{} {
	return []interface{}{}
}

func (c CircuitBreakerOption) Type() OptionType {
	return OptionTypeCircuitBreakOption
}
