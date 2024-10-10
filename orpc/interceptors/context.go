package interceptors

import (
	"context"

	"github.com/orion/orpc/xcontext"
	"github.com/orion/pkg/logger"
	"google.golang.org/grpc"
)

func ContextWrapperInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = xcontext.WrapContext(ctx)
		return handler(ctx, req)
	}
}
