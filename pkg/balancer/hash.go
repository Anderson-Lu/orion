package balancer

import "hash/crc32"

func NewHashBalancer(size int) *HashBalancer {
	return &HashBalancer{size: size}
}

type HashBalancer struct {
	size int
}

func (e *HashBalancer) Get(params ...interface{}) int {
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

func (e *HashBalancer) Update(index int) {}
