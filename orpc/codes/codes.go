package codes

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrCodeNone        = 0    // none
	ErrCodeUndefined   = 4000 // some errors not defined
	ErrCodeRateLimited = 4001 // ratelimit
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
