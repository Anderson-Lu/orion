package interceptors

import (
	"context"
	"errors"
	"fmt"

	"github.com/Anderson-Lu/orion/orpc/codes"
	"github.com/Anderson-Lu/orion/pkg/logger"
	"github.com/Anderson-Lu/orion/pkg/ratelimit"
	"google.golang.org/grpc"
)

var (
	lims = ratelimit.NewRateLimiters()
)

func RateLimitorInterceptor(configs []*ratelimit.Config, lg *logger.Logger) grpc.UnaryServerInterceptor {
	for _, v := range configs {
		lims.Register(v)
		lg.Info("[limitor] limitor registed", "key", v.Key, "cap", v.Cap, "tokensPerSec", v.TokensPerSecond)
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := info.FullMethod

		if !lims.Allow(method) {
			return nil, codes.WrapCodeFromError(errors.New("rate limited"), codes.ErrCodeRateLimited)
		}
		fmt.Println("--->", method)
		h, err := handler(ctx, req)
		return h, err
	}
}
