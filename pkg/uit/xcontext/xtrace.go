package xcontext

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const (
	KEY_X_FORWARDED_FOR = "x-forwarded-for"
	KEY_X_REQUEST_ID    = "x-request-id"
)

type TracerHeader struct {
	RequestId string
	Forward   string
	Uid       string
}

func BuildTraceHeader(ctx context.Context) *TracerHeader {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &TracerHeader{}
	}
	fmt.Println(md)
	return &TracerHeader{
		RequestId: getFirstElementFromArray(md.Get(KEY_X_REQUEST_ID)),
		Forward:   getFirstElementFromArray(md.Get(KEY_X_FORWARDED_FOR)),
	}
}

func getFirstElementFromArray(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[0]
}
