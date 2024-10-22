package options

import (
	"github.com/Anderson-Lu/orion/orpc/codec"
	"google.golang.org/grpc"
)

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
