package options

type OptionType int

const (
	OptionTypeGrpcCall     OptionType = 1
	OptionTypeBalancer     OptionType = 2
	OptionTypeCircuitBreak OptionType = 3
	OptionTypeResovler     OptionType = 4
	OptionTypeMethod       OptionType = 5
	OptionTypeMetadata     OptionType = 6
)

type OrionClientInvokeOption interface {
	Params() []interface{}
	Type() OptionType
}
