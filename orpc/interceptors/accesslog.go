package interceptors

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/orion/orpc/xcontext"
	"github.com/orion/pkg/logger"
	"google.golang.org/grpc"
)

func ChainInterceptors(ics ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return grpc_middleware.ChainUnaryServer(ics...)
}

func AccessInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := info.FullMethod
		begin := time.Now()
		h, err := handler(ctx, req)
		if lg == nil {
			return h, err
		}
		header := xcontext.BuildTraceHeader(ctx)
		if err == nil {
			lg.Info("[succ]", "method", method, "requestId", header.RequestId, "forward", header.Forward, "req", req, "rsp", h, "cost", time.Since(begin).Milliseconds())
		} else {
			lg.Error("[fail]", "method", method, "requestId", header.RequestId, "forward", header.Forward, "req", req, "rsp", h, "err", err, "cost", time.Since(begin).Milliseconds())
		}
		return h, err
	}
}
