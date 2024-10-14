package balancer

import (
	"math/rand"
	"time"
)

func NewRandomBalancer(size int) *RandomBalancer {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &RandomBalancer{size: size, r: r}
}

type RandomBalancer struct {
	size int
	r    *rand.Rand
}

func (e *RandomBalancer) Get(params ...interface{}) int {
	if e.size <= 0 {
		return 0
	}
	return e.r.Intn(e.size)
}

func (e *RandomBalancer) Update(index int) {}
