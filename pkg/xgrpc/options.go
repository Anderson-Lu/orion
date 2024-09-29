package xgrpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ServerOption func(s *Server)

func WithGRPCHandler(handler interface{}, sds ...*grpc.ServiceDesc) ServerOption {
	return func(s *Server) {
		if s.gServer == nil {
			return
		}
		for _, sd := range sds {
			s.gServer.RegisterService(sd, handler)
		}
	}
}

func WithGrpcGatewayEndpointFunc(rFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) ServerOption {
	return func(s *Server) {
		s.gatewayFunc = rFunc
	}
}

func WithFlags() ServerOption {
	return func(s *Server) {
		s.cmdMode = true
	}
}
