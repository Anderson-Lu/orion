package backoff

import (
	"sync"
	"sync/atomic"
)

type IBackOff interface {
	Next() (bool, int64)
	Reset()
}

var (
	_ IBackOff = (*ExponentialBackOff)(nil)
	_ IBackOff = (*LinearBackoff)(nil)
)

type ExponentialBackOff struct {
	max int64
	cur int64

	mu sync.Mutex
}

func (e *ExponentialBackOff) Next() (bool, int64) {

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.cur <= 0 {
		e.cur = 1
		return true, e.cur
	}

	e.cur++

	if e.cur > e.max {
		return false, 0
	}

	return true, 1 << (e.cur - 1)
}

func (e *ExponentialBackOff) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cur = 1
}

func NewExponentialBackOff(max int64) *LinearBackoff {
	return &LinearBackoff{
		max: max,
	}
}

func NewLinearBackoff(gap int64) *LinearBackoff {
	return &LinearBackoff{c: gap}
}

type LinearBackoff struct {
	c   int64
	max int64
	cur int64
}

func (e *LinearBackoff) Next() (bool, int64) {
	n := atomic.AddInt64(&e.cur, 1)
	if n > e.max {
		return false, 0
	}
	return true, e.c
}

func (e *LinearBackoff) Reset() {
	atomic.StoreInt64(&e.cur, 0)
}
