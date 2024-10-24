package options

func WithCircuitBreak(params ...string) OrionClientInvokeOption {
	if len(params) == 0 {
		return &CallOptionWithCircuitBreaker{}
	}
	return &CallOptionWithCircuitBreaker{key: params[0]}
}

type CallOptionWithCircuitBreaker struct {
	key string
}

func (c CallOptionWithCircuitBreaker) Params() []interface{} {
	return []interface{}{}
}

func (c CallOptionWithCircuitBreaker) Type() OptionType {
	return OptionTypeCircuitBreakOption
}

func (c CallOptionWithCircuitBreaker) Key() string {
	return c.key
}
