package xgrpc

import (
	"context"
	"time"

	"github.com/uit/pkg/logger"
	"google.golang.org/grpc"
)

func AccessLoggerInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := info.FullMethod
		begin := time.Now()
		h, err := handler(ctx, req)
		if lg == nil {
			return h, err
		}
		if err == nil {
			lg.Info("[succ]", "method", method, "cost", time.Since(begin).Milliseconds(), "req", req, "rsp", h)
		} else {
			lg.Error("[fail]", "method", method, "cost", time.Since(begin).Milliseconds(), "req", req, "rsp", h, "err", err)
		}
		return h, err
	}
}
