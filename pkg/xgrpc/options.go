package xgrpc

import (
	"google.golang.org/grpc"
)

type ServerOption func(s *Server)

func WithGRPCHandler(handler interface{}, sds ...*grpc.ServiceDesc) ServerOption {
	return func(s *Server) {
		if s.g == nil {
			return
		}
		for _, sd := range sds {
			s.g.RegisterService(sd, handler)
		}
	}
}

func WithFlags() ServerOption {
	return func(s *Server) {
		s.cmdMode = true
	}
}
