package options

import "google.golang.org/grpc/metadata"

func WithHeaders(kvs ...string) OrionClientInvokeOption {
	if len(kvs)%2 != 0 {
		return &CallOptionWithHeader{headers: map[string]string{}}
	}
	m := make(map[string]string)
	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return &CallOptionWithHeader{headers: m}
}

type CallOptionWithHeader struct {
	headers map[string]string
}

func (c CallOptionWithHeader) Params() []interface{} {
	return []interface{}{}
}

func (c CallOptionWithHeader) Type() OptionType {
	return OptionTypeMetadata
}

func (h CallOptionWithHeader) Metadata() metadata.MD {
	if h.headers == nil {
		return nil
	}
	hds := metadata.MD{}
	for k, v := range h.headers {
		hds.Set(k, v)
	}
	return hds
}
