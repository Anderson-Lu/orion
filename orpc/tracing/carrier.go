package tracing

import "google.golang.org/grpc/metadata"

type MetadataCarrier struct {
	md metadata.MD
}

func NewMetadataCarrier(md metadata.MD) MetadataCarrier {
	return MetadataCarrier{md: md}
}

func (m MetadataCarrier) Get(key string) string {
	if c := m.md.Get(key); len(c) > 0 {
		return c[0]
	}
	return ""
}

func (m MetadataCarrier) Set(key string, value string) {
	m.md.Set(key, value)
}

func (m MetadataCarrier) Keys() []string {
	r := []string{}
	for k := range m.md {
		r = append(r, k)
	}
	return r
}

func (m MetadataCarrier) GetTracerName() string {
	return m.Get(headerKeyTracerName)
}

func (m MetadataCarrier) SetTracerNameNx(tracerName string) string {
	if c := m.Get(headerKeyTracerName); c != "" {
		return c
	}
	m.Set(headerKeyTracerName, tracerName)
	return tracerName
}
