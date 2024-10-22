package options

func WithHeaders(kvs ...string) OrionClientInvokeOption {
	if len(kvs)%2 != 0 {
		return &HeaderOption{headers: map[string]string{}}
	}
	m := make(map[string]string)
	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return &HeaderOption{headers: m}
}

type HeaderOption struct {
	headers map[string]string
}

func (c HeaderOption) Params() []interface{} {
	return []interface{}{}
}

func (c HeaderOption) Type() OptionType {
	return OptionTypeCircuitBreakOption
}
