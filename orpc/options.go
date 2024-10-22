package orpc

import (
	"context"
	"errors"

	"github.com/Anderson-Lu/orion/orpc/parser"
	"github.com/Anderson-Lu/orion/orpc/registry"
	"github.com/Anderson-Lu/orion/orpc/registry/orion_consul"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type ServerOption func(s *Server) error

func WithConfig(c *Config) ServerOption {
	return func(s *Server) error {
		s.c = c
		return nil
	}
}

func WithConfigFile(fsPath string) ServerOption {
	return func(s *Server) error {
		var cc = &Config{}

		if err := parser.ParseConfigFile(fsPath, cc); err != nil {
			return errors.New("parse config file error:" + err.Error())
		}

		s.c = cc
		return nil
	}
}

func WithGRPCHandler(handler interface{}, sds ...*grpc.ServiceDesc) ServerOption {
	return func(s *Server) error {
		if s.grpcHandlers == nil {
			s.grpcHandlers = make(map[interface{}][]*grpc.ServiceDesc)
		}
		s.grpcHandlers[handler] = sds
		return nil
	}
}

func WithGrpcGatewayEndpointFunc(rFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) ServerOption {
	return func(s *Server) error {
		s.gatewayFunc = rFunc
		return nil
	}
}

func WithFlags() ServerOption {
	return func(s *Server) error {
		s.cmdMode = true
		return nil
	}
}

func WithRegistry(r registry.RegisteyMode) ServerOption {
	return func(s *Server) error {
		switch r {
		case registry.RegisteyConsul:
			if s.c == nil || s.c.Registry == nil {
				return errors.New("nil config, please configure first")
			}
			s.rsy = orion_consul.NewOrionConsulRegistry(s.c.Registry.Service, s.c.Registry.IP, s.c.Registry.Port)
		}
		return nil
	}
}
