package resolver

import "google.golang.org/grpc"

type IResolver interface {
	Name() string
	Select(serviceName string, params ...interface{}) (*grpc.ClientConn, error)
	SetDialOptions(opts ...grpc.DialOption)
}

var (
	_ IResolver = (*ConsulResovler)(nil)
	_ IResolver = (*DirectResolver)(nil)
)
