package orion_consul

import (
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type OrionConsulWatchers struct {
	address string
	opts    []interface{}
	ws      map[string]*OrionConsulWatcher
	mu      sync.RWMutex
}

func NewOrionConsulWatchers(address string, opts ...interface{}) *OrionConsulWatchers {
	o := &OrionConsulWatchers{address: address, opts: opts, ws: make(map[string]*OrionConsulWatcher)}
	return o
}

func (o *OrionConsulWatchers) Watch(service string) error {

	precheck := func() bool {
		o.mu.RLock()
		defer o.mu.RUnlock()
		if _, ok := o.ws[service]; ok {
			return true
		}
		return false
	}

	if precheck() {
		return nil
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	w, err := NewOrionConsulWatcher(o.address, service, o.opts...)
	if err != nil {
		return err
	}

	o.ws[service] = w
	w.Run()
	return nil
}

type OrionConsulWatcher struct {
	service    string
	address    string
	p          *watch.Plan
	notifyFunc OrionNodesNotifyFunc
}

func NewOrionConsulWatcher(address, service string, opts ...interface{}) (*OrionConsulWatcher, error) {
	o := &OrionConsulWatcher{service: service, address: address}

	for _, v := range opts {
		if n, ok := v.(OrionNodesNotifyFunc); ok {
			o.notifyFunc = n
		}
	}

	if err := o.init(); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *OrionConsulWatcher) init() error {
	p, err := watch.Parse(map[string]interface{}{"type": "service", "service": o.service})
	if err != nil {
		return err
	}
	o.p = p
	p.Handler = func(u uint64, i interface{}) {
		switch ii := i.(type) {
		case []*api.ServiceEntry:
			nodes := []OrionNode{}
			for _, v := range ii {
				nodes = append(nodes, OrionNode{
					Namespace:   v.Service.Namespace,
					Datacenter:  v.Service.Datacenter,
					Host:        v.Service.Address,
					Port:        v.Service.Port,
					ServiceName: v.Service.Service,
					Tags:        v.Service.Tags,
				})
				if o.notifyFunc != nil {
					o.notifyFunc(o.service, nodes)
				}
			}
		}
	}
	return nil
}

func (o *OrionConsulWatcher) Run() {
	if o.p == nil {
		return
	}
	o.p.Run(o.address)
}
