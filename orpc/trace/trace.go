package trace

type TraceHeader struct{}

type TraceBody interface{}

type Trace struct {
	TraceHeader
	TraceBody
}

type Tracer struct{}

func (t *Tracer) Report() {}
