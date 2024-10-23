package resolver

import (
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewDirectResolver(address string) IResolver {
	return &DirectResolver{address: address}
}

type DirectResolver struct {
	c       *grpc.ClientConn
	address string
	sg      singleflight.Group
	inited  bool
}

func (d *DirectResolver) Name() string {
	return "default"
}

func (d *DirectResolver) Select(serviceName string) (*grpc.ClientConn, error) {
	c, err, _ := d.sg.Do("init", func() (interface{}, error) {
		if d.inited {
			return d.c, nil
		}
		opts := []grpc.DialOption{}
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		c, err := grpc.NewClient(d.address, opts...)
		if err != nil {
			return nil, err
		}
		d.c = c
		d.inited = true
		return d.c, nil
	})
	if err != nil {
		return nil, err
	}
	return c.(*grpc.ClientConn), nil
}
