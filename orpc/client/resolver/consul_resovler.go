package resolver

import (
	"fmt"
	"sync"

	"github.com/Anderson-Lu/orion/orpc/registry/orion_consul"
	"google.golang.org/grpc"
)

type ConsulResovler struct {
	address  string
	conns    map[string]*grpc.ClientConn
	mu       sync.RWMutex
	watchers *orion_consul.OrionConsulWatchers
}

func NewConsulResovler(address string) *ConsulResovler {
	c := &ConsulResovler{address: address, conns: make(map[string]*grpc.ClientConn)}
	c.watchers = orion_consul.NewOrionConsulWatchers(c.address, c.notify)
	return c
}

func (c *ConsulResovler) key(service string, address string) string {
	return fmt.Sprintf("conn_%s_%s", service, address)
}

func (c *ConsulResovler) run() {

}

func (c *ConsulResovler) notify(nodes []orion_consul.OrionNode) {
	for _, v := range nodes {
		key := c.key(v.ServiceName, fmt.Sprintf("%s:%d", v.Host, v.Port))
	}
}

func (d *DirectResolver) Select(serviceName string) (*grpc.ClientConn, error) {

}
