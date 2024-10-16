package main

import (
	"fmt"
	"time"

	"github.com/Anderson-Lu/orion/pkg/circuit_break"
)

func main() {

	rc := &circuit_break.RuleConfig{
		Name: "rule1",
		Window: &circuit_break.WindowConfig{
			Duration: 10000,
			Buckets:  10,
		},
		OpenDuration:     1000,
		HalfOpenDuration: 1000,
		HaflOpenPassRate: 0,
		RuleExpression:   "req_count >= 1 && succ_rate < 0.90",
	}

	cb := circuit_break.NewCircuitBreaker()
	cb.Register(rc)

	for i := 0; i < 2000; i++ {

		time.Sleep(time.Millisecond * 800)
		if pass := cb.Pass(rc.Name); pass {
			fmt.Println("正常请求了...")
		} else {
			fmt.Println("被拦截了...")
			continue
		}
		// dosomething

		mockResult := false

		if i > 50 {
			mockResult = true
		}
		cb.Report(rc.Name, mockResult, int64(i%10))
	}
}
