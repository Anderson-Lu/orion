package options

import "google.golang.org/grpc"

type BaseServer interface {
	RegisterGRPCHandler(sd *grpc.ServiceDesc, handler interface{})
	RegisterFlagsHandler()
	ListenAndServe() error
}

type ServerRoute struct {
	Handler interface{}
	Sds     []*grpc.ServiceDesc
}

type ServerOption func(s BaseServer)

func WithHandler(handler interface{}, sds ...*grpc.ServiceDesc) ServerOption {
	return func(s BaseServer) {
		for _, sd := range sds {
			s.RegisterGRPCHandler(sd, handler)
		}
	}
}

func WithFlags() ServerOption {
	return func(s BaseServer) {
		s.RegisterFlagsHandler()
	}
}
