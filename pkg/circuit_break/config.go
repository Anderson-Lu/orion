package circuit_break

type RuleConfig struct {

	// the name of the rule.
	Name string

	// specify the stat window including buckets size and duration.
	Window *WindowConfig

	OpenDuration     int64 // if circuit-break occours, it will recover after next BreakDuration duration
	HalfOpenDuration int64 // eg: 60s = 60 * 1000 = 60000
	HaflOpenPassRate int64 // 0-100

	// specify different rule checker to handle circuit break.
	BreakHandlerType string

	// expression for checking,
	// suuport ops: [ +, -, *, /, &&, || ]
	// support vars: [req_count, req_succ_count, req_fail_count, succ_rate]
	// eg: `req_count > 100 && succ_rate < 0.95`
	RuleExpression string
}
