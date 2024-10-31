package interceptors

import (
	"context"

	"github.com/Anderson-Lu/orion/orpc/tracing"
	"github.com/Anderson-Lu/orion/pkg/logger"
	"go.opentelemetry.io/otel"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ContextWrapperInterceptor(lg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		carrier := tracing.NewMetadataCarrier(md)
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		return handler(ctx, req)
	}
}
