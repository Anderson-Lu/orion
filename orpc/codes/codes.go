package codes

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrCodeNone                     = 0    // none
	ErrCodeCircuitBreak             = 3001 // circuit breaked
	ErrCodeClientConnNotEstablished = 3100 // client grpc connection not established
	ErrCodeUndefined                = 4000 // some errors not defined
	ErrCodeRateLimited              = 4001 // ratelimit
)

var (
	ErrClientConnNotEstablished = WrapCodeFromError(errors.New("client conn not established yet"), ErrCodeClientConnNotEstablished)
	ErrClientCircuitBreaked     = WrapCodeFromError(errors.New("circuit broken"), ErrCodeCircuitBreak)
)

func GetCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	st, ok := status.FromError(err)
	if !ok {
		return ErrCodeUndefined
	}
	return int(st.Code())
}

func GetCodeAndMessageFromError(err error) (int, string) {
	if err == nil {
		return 0, ""
	}
	st, ok := status.FromError(err)
	if !ok {
		return ErrCodeUndefined, "undefined error:" + err.Error()
	}
	return int(st.Code()), st.Message()
}

func WrapCodeFromError(err error, code int) error {
	if err == nil {
		return nil
	}
	return status.Error(codes.Code(code), err.Error())
}
