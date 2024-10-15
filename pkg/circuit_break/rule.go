package circuit_break

import (
	"context"
	"time"
)

type IRuleChecker interface {
	// metrics is [[succCnt, failCnt, avgCost]]
	Turn(ctx context.Context, metrics [][]float64) CircuitBreakStatus
	Pass() bool
}

type Rule struct {
	rc *RuleConfig
	rr IRuleChecker
	md *Window
}

func NewRule(rc *RuleConfig, windowConfig *WindowConfig) *Rule {
	r := &Rule{}
	r.rc = rc
	switch rc.BreakHandlerType {
	default:
		r.rr = NewDefaultRuleStratege(rc.RuleExpression, rc.OpenDuration, rc.HalfOpenDuration, rc.HaflOpenPassRate)
	}
	r.md = NewWindow(windowConfig, r.rr)
	return r
}

func (r *Rule) Pass() bool {
	return r.rr.Pass()
}

func (r *Rule) Report(result bool, cost int64) {
	if result {
		r.md.Succ(time.Now().UnixMilli(), cost)
		return
	}
	r.md.Fail(time.Now().UnixMilli(), cost)
}
