package interceptors

import (
	"runtime/debug"

	"github.com/Anderson-Lu/orion/pkg/logger"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		lg.Error("[panic]", "panicInfo", p, "stack", string(debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}))
}
