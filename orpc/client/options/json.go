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
	return []interface{}{}
}

func (c CallOptionWithJson) Type() OptionType {
	return OptionTypeGrpcCall
}

func (c CallOptionWithJson) GrpcCallOption() grpc.CallOption {
	return grpc.CallContentSubtype(codec.JSON{}.Name())
}
