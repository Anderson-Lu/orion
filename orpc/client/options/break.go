package options

import "github.com/Anderson-Lu/orion/pkg/circuit_break"

func WithCircuitBreak(rules ...*circuit_break.RuleConfig) OrionClientInvokeOption {
	return &CircuitBreakerOption{rules: rules}
}

type CircuitBreakerOption struct {
	rules []*circuit_break.RuleConfig
}

func (c CircuitBreakerOption) Params() []interface{} {
	return []interface{}{}
}

func (c CircuitBreakerOption) Type() OptionType {
	return OptionTypeCircuitBreakOption
}

func (c CircuitBreakerOption) IsMatch(medhod string) (bool, *circuit_break.RuleConfig) {
	for _, v := range c.rules {
		if v.Name == medhod {
			return true, v
		}
	}
	return false, nil
}
