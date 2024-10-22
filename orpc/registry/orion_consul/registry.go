package orion_consul

import (
	"context"
	"fmt"

	"github.com/Anderson-Lu/orion/orpc/registry"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type OrionConsulRegistry struct {
	consuleIP   string
	consulePort uint32
	service     string
	ci          *api.Client
	cb          *health.Server
}

func NewOrionConsulRegistry(serviceName string, consuleIP string, consulPort uint32) *OrionConsulRegistry {
	rsy := &OrionConsulRegistry{}
	rsy.consuleIP = consuleIP
	rsy.consulePort = consulPort
	rsy.service = serviceName
	return rsy
}

func (o *OrionConsulRegistry) init() error {

	if o.ci != nil {
		return nil
	}

	c := api.DefaultConfig()
	c.Address = fmt.Sprintf("%s:%d", o.consuleIP, o.consulePort)
	ci, err := api.NewClient(c)
	if err != nil {
		return err
	}
	o.ci = ci
	return nil
}

func (o *OrionConsulRegistry) RegisterHealthHandler(svr *grpc.Server) {
	o.cb = health.NewServer()
	grpc_health_v1.RegisterHealthServer(svr, o.cb)
}

func (o *OrionConsulRegistry) AddNode(ctx context.Context, ip string, port uint32, opts ...registry.IRegistryAddOption) error {

	if err := o.init(); err != nil {
		return err
	}

	if o.service == "" {
		o.service = fmt.Sprintf("service:%s:%d", ip, port)
	}

	node := new(api.AgentServiceRegistration)
	node.ID = o.service
	node.Name = o.service
	node.Port = int(port)
	node.Address = ip
	node.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, port),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}

	for _, v := range opts {
		switch opt := v.(type) {
		case *registry.RegistryAddOptionTags:
			node.Tags = append(node.Tags, (*opt)...)
		}
	}

	return o.ci.Agent().ServiceRegister(node)
}

func (o *OrionConsulRegistry) GetNodes(ctx context.Context, name string) ([]*registry.Node, error) {

	serviceMap, err := o.ci.Agent().ServicesWithFilter(fmt.Sprintf("Service==`%s`", name))
	if err != nil {
		return nil, err
	}

	r := []*registry.Node{}
	for _, v := range serviceMap {
		r = append(r, &registry.Node{
			IP:   v.Address,
			Port: uint32(v.Port),
		})
	}

	return r, nil
}

func (o *OrionConsulRegistry) RemoveNode(ctx context.Context) error {
	o.ci.Agent().ServiceDeregister(o.service)
	return nil
}
