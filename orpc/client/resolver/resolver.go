package resolver

import "google.golang.org/grpc"

type IResolver interface {
	Select() (*grpc.ClientConn, error)
}
