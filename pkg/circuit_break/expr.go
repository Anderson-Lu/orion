package circuit_break

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

const (
	FieldRequestNum     string = "req_count"
	FieldRequestSuccNum string = "req_succ_count"
	FieldRequestFailNum string = "req_fail_count"
	FieldSuccRate       string = "succ_rate"
	FieldAvgCost        string = "avg_cost"
)

type CircuitRuleAst struct {
	ruleString string
	expr       ast.Expr
	exprErr    error
}

func (c *CircuitRuleAst) parse(params map[string]float64) (bool, error) {

	if c.exprErr != nil {
		return false, c.exprErr
	}

	if c.expr == nil {
		expr, err := parser.ParseExpr(c.ruleString)
		if err != nil {
			c.exprErr = err
			return false, err
		}
		c.expr = expr
	}

	switch c.expr.(type) {
	case *ast.BinaryExpr:
		result := c.eval(c.expr, params)
		return result == 1, nil
	default:
		return false, errors.New("invalid expr found, only binary-expression supported")
	}
}

func (c *CircuitRuleAst) eval(exp ast.Expr, vars map[string]float64) float64 {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return c.evalBinary(exp, vars)
	case *ast.BasicLit:
		f, _ := strconv.ParseFloat(exp.Value, 64)
		return f
	case *ast.Ident:
		return vars[exp.Name]
	}
	return 0
}

func (c *CircuitRuleAst) bool2float(ok bool) float64 {
	if ok {
		return 1
	}
	return 0
}

func (c *CircuitRuleAst) evalBinary(exp *ast.BinaryExpr, vars map[string]float64) float64 {
	switch exp.Op {
	case token.ADD:
		return c.eval(exp.X, vars) + c.eval(exp.Y, vars)
	case token.MUL:
		return c.eval(exp.X, vars) * c.eval(exp.Y, vars)
	case token.SUB:
		return c.eval(exp.X, vars) - c.eval(exp.Y, vars)
	case token.QUO:
		return c.eval(exp.X, vars) / c.eval(exp.Y, vars)
	case token.GTR:
		return c.bool2float(c.eval(exp.X, vars) > c.eval(exp.Y, vars))
	case token.GEQ:
		return c.bool2float(c.eval(exp.X, vars) >= c.eval(exp.Y, vars))
	case token.EQL:
		return c.bool2float(c.eval(exp.X, vars) == c.eval(exp.Y, vars))
	case token.LSS:
		return c.bool2float(c.eval(exp.X, vars) < c.eval(exp.Y, vars))
	case token.LEQ:
		return c.bool2float(c.eval(exp.X, vars) <= c.eval(exp.Y, vars))
	case token.LAND:
		return c.bool2float(c.eval(exp.X, vars) > 0 && c.eval(exp.Y, vars) > 0)
	case token.LOR:
		return c.bool2float(c.eval(exp.X, vars) > 0 || c.eval(exp.Y, vars) > 0)
	}
	return 0
}
