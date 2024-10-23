package circuit_break

import "sync"

type CircuitBreakStatus int

const (
	CircuitBreakStatusClose    CircuitBreakStatus = 0
	CircuitBreakStatusHalfOpen CircuitBreakStatus = 1
	CircuitBreakStatusOpen     CircuitBreakStatus = 2
)

type CircuitBreaker struct {
	mu    sync.RWMutex
	rules map[string]*Rule
}

func NewCircuitBreaker() *CircuitBreaker {
	cb := &CircuitBreaker{}
	cb.rules = make(map[string]*Rule)
	return cb
}

func (c *CircuitBreaker) Register(rc *RuleConfig) {
	if c.rules == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.rules[rc.Name]; ok {
		return
	}
	c.rules[rc.Name] = NewRule(rc, rc.Window)
}

func (c *CircuitBreaker) Pass(resourceId string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r := c.rules[resourceId]
	if r == nil {
		return true
	}
	return r.Pass()
}

func (c *CircuitBreaker) Report(resourceId string, result bool, cost int64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r := c.rules[resourceId]
	if r == nil {
		return
	}
	r.Report(result, cost)
}
