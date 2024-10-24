package options

func WithBalancerParams(params ...interface{}) OrionClientInvokeOption {
	if len(params) == 0 {
		return &CallOptionWithBalancerParams{}
	}
	return &CallOptionWithBalancerParams{params: params}
}

type CallOptionWithBalancerParams struct {
	params []interface{}
}

func (c CallOptionWithBalancerParams) Params() []interface{} {
	return c.params
}

func (c CallOptionWithBalancerParams) Type() OptionType {
	return OptionTypeBalancer
}
