package tracing

type TracingOption func(*Tracing)

func WithOpenTelemetryAddress(address string) TracingOption {
	return func(t *Tracing) {
		t.exportorAddr = address
	}
}

func WithResource(r *Resources) TracingOption {
	return func(t *Tracing) {
		t.baseResources = r
	}
}
