package interceptors

import (
	"context"
	"time"

	"github.com/uit/pkg/logger"
	"google.golang.org/grpc"
)

func AccessInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := info.FullMethod
		begin := time.Now()
		h, err := handler(ctx, req)
		if lg == nil {
			return h, err
		}
		tick := time.Now()
		if err == nil {
			lg.Info("[succ]", "method", method, "cost", time.Since(begin).Milliseconds(), "req", req, "rsp", h, "cost", time.Since(tick).Milliseconds())
		} else {
			lg.Error("[fail]", "method", method, "cost", time.Since(begin).Milliseconds(), "req", req, "rsp", h, "err", err, "cost", time.Since(tick).Milliseconds())
		}
		return h, err
	}
}
