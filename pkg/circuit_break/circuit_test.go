package circuit_break

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreak(t *testing.T) {

	rc := &RuleConfig{
		Name: "rule1",
		Window: &WindowConfig{
			Duration: 10000,
			Buckets:  10,
		},
		OpenDuration:     1000,
		HalfOpenDuration: 1000,
		HaflOpenPassRate: 0,
		RuleExpression:   "req_count > 10 && succ_rate < 0.5",
	}

	cb := NewCircuitBreaker()
	cb.Register(rc)

	assert.True(t, cb.Pass(rc.Name), "must be true on init status")
	for i := 0; i < 10; i++ {
		cb.Report(rc.Name, true, 10)
	}
	assert.True(t, cb.Pass(rc.Name), "must be true on init status")
	time.Sleep(time.Second)
	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 10)
		cb.Report(rc.Name, false, 10)
	}
	time.Sleep(time.Second)
	fmt.Println(cb.Pass(rc.Name), "+==")
	assert.False(t, cb.Pass(rc.Name), "must be true on init status")
}
