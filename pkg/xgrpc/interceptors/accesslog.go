package interceptors

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/uit/pkg/logger"
	"github.com/uit/pkg/xgrpc/xcontext"
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
			lg.Info("[succ]", "method", method, "cost", "request-id", header.RequestId, "clientIP", header.ClientIP, time.Since(begin).Milliseconds(), "req", req, "rsp", h)
		} else {
			lg.Error("[fail]", "method", method, "cost", "request-id", header.RequestId, "clientIP", header.ClientIP, time.Since(begin).Milliseconds(), "req", req, "rsp", h, "err", err)
		}
		return h, err
	}
}
