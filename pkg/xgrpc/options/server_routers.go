package options

import "google.golang.org/grpc"

type BaseServer interface {
	Register(sd *grpc.ServiceDesc, handler interface{})
}

type ServerRoute struct {
	Handler interface{}
	Sds     []*grpc.ServiceDesc
}

type ServerOption func(s BaseServer)

func WithHandler(handler interface{}, sds ...*grpc.ServiceDesc) ServerOption {
	return func(s BaseServer) {
		for _, sd := range sds {
			s.Register(sd, handler)
		}
	}
}
