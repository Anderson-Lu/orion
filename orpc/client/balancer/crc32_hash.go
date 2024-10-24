package balancer

import (
	"hash/crc32"
)

func NewCrc32HashBalancer() *Crc32HashBalancer {
	return &Crc32HashBalancer{}
}

type Crc32HashBalancer struct {
	size int
}

func (e *Crc32HashBalancer) Get(params ...interface{}) int {
	if e.size <= 1 {
		return 0
	}
	if len(params) <= 0 {
		return 0
	}
	if str, ok := params[1].(string); !ok {
		return 0
	} else {
		return int(crc32.ChecksumIEEE([]byte(str))) % e.size
	}
}

func (e *Crc32HashBalancer) Update(index int) {}

func (e *Crc32HashBalancer) Resize(size int) {
	if size < 0 {
		return
	}
	e.size = size
}

func (e *Crc32HashBalancer) Copy() Balancer {
	return &DefaultBalancer{size: e.size}
}
