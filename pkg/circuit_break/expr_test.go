package circuit_break

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr(t *testing.T) {
	ast := &CircuitRuleAst{ruleString: "req_count > 10 && req_succ_count < 60, 22"}
	ok, err := ast.parse(map[string]float64{"req_succ_count": 50.1})
	assert.NotNil(t, err, "invalid expr")
	assert.Equal(t, false, ok, "not ok")

	ast2 := &CircuitRuleAst{ruleString: "req_count > 10 && req_succ_count < 60"}
	ok, err = ast2.parse(map[string]float64{"req_count": 100, "req_succ_count": 50.1})
	assert.Nil(t, err, "invalid expr")
	assert.Equal(t, true, ok, "not ok")

	ok, err = ast2.parse(map[string]float64{"req_count": 1, "req_succ_count": 50.1})
	assert.Nil(t, err, "invalid expr")
	assert.Equal(t, false, ok, "not ok")
}
