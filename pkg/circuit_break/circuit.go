package circuit_break

type CircuitBreakStatus int

const (
	CircuitBreakStatusClose    CircuitBreakStatus = 0
	CircuitBreakStatusHalfOpen CircuitBreakStatus = 1
	CircuitBreakStatusOpen     CircuitBreakStatus = 2
)

type CircuitBreaker struct {
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
	if _, ok := c.rules[rc.Name]; ok {
		return
	}
	c.rules[rc.Name] = NewRule(rc, rc.Window)
}

func (c *CircuitBreaker) Pass(resourceId string) bool {
	r := c.rules[resourceId]
	if r == nil {
		return true
	}
	return r.Pass()
}

func (c *CircuitBreaker) Report(resourceId string, result bool, cost int64) {
	r := c.rules[resourceId]
	if r == nil {
		return
	}
	r.Report(result, cost)
}
