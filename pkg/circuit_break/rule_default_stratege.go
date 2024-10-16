package circuit_break

import (
	"context"
	"math/rand"
	"time"
)

func NewDefaultRuleStratege(exprString string, openDuration, halfDuration, halfPassRate int64) *RuleDefaultStratege {
	r := &RuleDefaultStratege{
		ast:          &CircuitRuleAst{ruleString: exprString},
		halfRand:     rand.New(rand.NewSource(time.Now().UnixMicro())),
		halfDuration: halfDuration,
		halfPassRate: halfPassRate,
		openDuration: openDuration,
	}
	r.ast.parse(nil)
	return r
}

type RuleDefaultStratege struct {
	ast *CircuitRuleAst

	openDuration int64

	halfDuration int64
	halfPassRate int64
	halfRand     *rand.Rand

	lastTurn   int64
	lastStatus CircuitBreakStatus
}

func (r *RuleDefaultStratege) Turn(ctx context.Context, metrics [][]float64) CircuitBreakStatus {

	if len(metrics) == 0 {
		return r.lastStatus
	}

	now := time.Now().UnixMilli()

	switch r.lastStatus {

	case CircuitBreakStatusOpen:
		if now-r.lastTurn < r.openDuration {
			return r.lastStatus
		}
		if r.ok(metrics) {
			r.lastStatus = CircuitBreakStatusHalfOpen
		}
		r.lastTurn = now
		return r.lastStatus

	case CircuitBreakStatusHalfOpen:
		if time.Now().UnixMilli()-r.lastTurn < r.halfDuration {
			return r.lastStatus
		}

		if r.ok(metrics) {
			r.lastStatus = CircuitBreakStatusClose
		} else {
			r.lastStatus = CircuitBreakStatusOpen
		}
		r.lastTurn = now
		return r.lastStatus

	case CircuitBreakStatusClose:
		if !r.ok(metrics) {
			r.lastStatus = CircuitBreakStatusHalfOpen
		}
		r.lastTurn = now
		return r.lastStatus
	}

	return r.lastStatus
}

func (r *RuleDefaultStratege) ok(metrics [][]float64) bool {

	var cnt, succ, fail, avgCost float64
	for _, labels := range metrics {
		cnt += (labels[0] + labels[1])
		succ += labels[0]
		fail += labels[1]
		avgCost += labels[2]
	}
	avgCost = avgCost / float64(len(metrics))
	succRate := 0.0
	if succ > 0 {
		succRate = cnt / succ
	}

	if cnt == 0 {
		return true
	}

	matched, err := r.ast.parse(map[string]float64{
		FieldRequestNum:     cnt,
		FieldRequestSuccNum: succ,
		FieldRequestFailNum: fail,
		FieldSuccRate:       succRate,
		FieldAvgCost:        avgCost,
	})
	if err != nil {
		return false
	}

	return !matched
}

func (r *RuleDefaultStratege) Pass() bool {
	switch r.lastStatus {
	case CircuitBreakStatusClose:
		return true
	case CircuitBreakStatusHalfOpen:
		return r.halfPass()
	case CircuitBreakStatusOpen:
		return false
	}
	return false
}

func (r *RuleDefaultStratege) halfPass() bool {
	if r.halfPassRate <= 0 {
		return false
	}
	if r.halfPassRate >= 100 {
		return true
	}
	return r.halfRand.Intn(100) > int(r.halfPassRate)
}
