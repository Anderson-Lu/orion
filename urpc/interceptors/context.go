package interceptors

import (
	"context"

	"github.com/uit/modules/logger"
	"github.com/uit/urpc/xcontext"
	"google.golang.org/grpc"
)

func ContextWrapperInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = xcontext.WrapContext(ctx)
		return handler(ctx, req)
	}
}