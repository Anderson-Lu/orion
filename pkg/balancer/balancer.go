package balancer

import "sync/atomic"

type Balancer interface {
	Get(params ...interface{}) int
	Update(index int)
}

func NewDefaultBalancer(size int) *DefaultBalancer {
	return &DefaultBalancer{
		size: size,
	}
}

type DefaultBalancer struct {
	size int
	c    int32
}

func (e *DefaultBalancer) Get(params ...interface{}) int {
	if e.size <= 1 {
		return 0
	}
	new := atomic.AddInt32(&e.c, 1)

	if new > int32(e.size)*2 {
		// it is not accurate under concurrent conditions, but it does not affect
		atomic.StoreInt32(&e.c, 0)
	}
	return int(new) % e.size
}

func (e *DefaultBalancer) Update(index int) {}
