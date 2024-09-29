package xcontext

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func WrapContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	return metadata.NewIncomingContext(ctx, filterTraceMD(md))
}
