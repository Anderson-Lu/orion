package resolver

import "google.golang.org/grpc"

type IResolver interface {
	Name() string
	Select(serviceName string, params ...interface{}) (*grpc.ClientConn, error)
}
