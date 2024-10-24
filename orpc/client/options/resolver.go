package options

func WithService(service string) OrionClientInvokeOption {
	c := &CallOptionWithService{service: service}
	return c
}

func WithDirectAddress(addr string) OrionClientInvokeOption {
	c := &CallOptionWithService{direct: addr}
	return c
}

type CallOptionWithService struct {
	service string
	direct  string
}

func (c CallOptionWithService) Params() []interface{} {
	return []interface{}{}
}

func (c CallOptionWithService) Type() OptionType {
	return OptionTypeResovler
}

func (c CallOptionWithService) Service() string {
	return c.service
}

func (c CallOptionWithService) Direct() string {
	return c.direct
}

func (c CallOptionWithService) IsDirect() bool {
	return c.direct != ""
}
