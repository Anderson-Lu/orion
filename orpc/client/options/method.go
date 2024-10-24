package options

func WithMethod(method string) *CallOptionWithPath {
	return &CallOptionWithPath{method: method}
}

type CallOptionWithPath struct {
	method string
}

func (c CallOptionWithPath) Params() []interface{} {
	return []interface{}{}
}

func (c CallOptionWithPath) Type() OptionType {
	return OptionTYpeMethod
}

func (c CallOptionWithPath) Method() string {
	return c.method
}
