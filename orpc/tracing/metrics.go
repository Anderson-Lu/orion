package tracing

type Metrics map[string]string

func (m Metrics) Set(key, value string) Metrics {
	m[key] = value
	return m
}
