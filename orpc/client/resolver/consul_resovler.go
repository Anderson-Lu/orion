package resolver

import (
	"fmt"
	"sync"

	"github.com/Anderson-Lu/orion/orpc/codes"
	"github.com/Anderson-Lu/orion/orpc/registry/orion_consul"
	"github.com/Anderson-Lu/orion/pkg/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConsulResovler struct {
	address  string
	conns    map[string]*ConsulResovlerConns
	mu       sync.RWMutex
	watchers *orion_consul.OrionConsulWatchers
	b        balancer.Balancer
}

type ConsulResovlerOption func(*ConsulResovler)

func WithBalancer(b balancer.Balancer) ConsulResovlerOption {
	return func(cr *ConsulResovler) {
		cr.b = b
	}
}

func NewConsulResovler(address string, opts ...ConsulResovlerOption) *ConsulResovler {
	c := &ConsulResovler{address: address, conns: make(map[string]*ConsulResovlerConns)}
	c.watchers = orion_consul.NewOrionConsulWatchers(c.address, c.notify)

	for _, opt := range opts {
		opt(c)
	}

	if c.b == nil {
		c.b = balancer.NewDefaultBalancer()
	}
	return c
}

func (c *ConsulResovler) notify(serviceName string, nodes []orion_consul.OrionNode) {
	addrs := []string{}
	for _, v := range nodes {
		key := fmt.Sprintf("%s:%d", v.Host, v.Port)
		addrs = append(addrs, key)
	}
	c.update(serviceName, addrs)
}

func (c *ConsulResovler) update(service string, addrs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.conns[service]; !ok {
		c.conns[service] = &ConsulResovlerConns{b: c.b}
	}
	c.conns[service].update(addrs)
}

func (c *ConsulResovler) getConns(serviceName string) *ConsulResovlerConns {
	c.mu.RLock()
	defer c.mu.RUnlock()
	k, ok := c.conns[serviceName]
	if !ok {
		return nil
	}
	return k
}

func (c *ConsulResovler) Select(serviceName string) (*grpc.ClientConn, error) {
	if conn := c.getConns(serviceName); conn == nil {
		go c.watchers.Watch(serviceName)
		return nil, codes.ErrClientConnNotEstablished
	} else {
		return conn.Select(serviceName)
	}
}

func (c *ConsulResovler) Name() string {
	return "consul"
}

type ConsulResovlerConns struct {
	conns []*ConsulResovlerConn
	mu    sync.RWMutex
	b     balancer.Balancer
}

func (c *ConsulResovlerConns) update(incomingAddrs []string) {

	fmt.Println("incomingAddrs", incomingAddrs)

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(incomingAddrs) == 0 {
		c.conns = c.conns[0:0]
		return
	}

	incoming := make(map[string]struct{}, 0)
	for _, v := range incomingAddrs {
		incoming[v] = struct{}{}
	}

	master := make(map[string]int, 0)
	for i, v := range c.conns {
		master[v.addr] = i
	}

	var newConns []*ConsulResovlerConn
	for _, v := range incomingAddrs {
		idx, ok := master[v]
		if !ok {
			c, err := newConsulResovlerConn(v)
			if err == nil {
				newConns = append(newConns, c)
			}
			continue
		}
		newConns = append(newConns, c.conns[idx])
	}

	c.conns = newConns
	c.b.Resize(len(c.conns))
}

func (c *ConsulResovlerConns) Select(serviceName string) (*grpc.ClientConn, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.conns) == 0 {
		return nil, codes.ErrClientConnNotEstablished
	}

	connIdx := c.b.Get()
	c.b.Update(connIdx)

	return c.conns[connIdx].c, nil
}

type ConsulResovlerConn struct {
	c    *grpc.ClientConn
	addr string
}

func newConsulResovlerConn(addr string) (*ConsulResovlerConn, error) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	c, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}
	return &ConsulResovlerConn{c: c, addr: addr}, nil
}