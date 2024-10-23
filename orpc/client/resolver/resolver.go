package resolver

import "google.golang.org/grpc"

type IResolver interface {
	Name() string
	Select(serviceName string) (*grpc.ClientConn, error)
}
