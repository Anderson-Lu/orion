package options

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
