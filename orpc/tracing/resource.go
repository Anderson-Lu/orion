package tracing

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	KEY_RESOURCE_SERVICE_NAME = "service.name"
	KEY_SPAN_KIND             = "span.kind"
	KEY_RESOURCE_INSTANCE_ID  = "orion.instance.id"
	KEY_RESOURCE_INSTANCE_IP  = "orion.instance.ip"
	KEY_SPAN_ERRCODE          = "orion.code"
	KEY_RESOURCE_ENV          = "orion.environment"
	KEY_RESOURCE_NAMESPACE    = "orion.namespace"
	KEY_UNI_TRACE_ID          = "orion.traceid"
)

type Resources struct {
	kvs []attribute.KeyValue
}

func (m *Resources) SetEnv(env string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_ENV,
		Value: attribute.StringValue(env),
	})
}

func (m *Resources) SetInstanceId(id string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_INSTANCE_ID,
		Value: attribute.StringValue(id),
	})
}

func (m *Resources) SetIP(ipAddr string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_INSTANCE_IP,
		Value: attribute.StringValue(ipAddr),
	})
}

func (m *Resources) SetNamespace(namespace string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_NAMESPACE,
		Value: attribute.StringValue(namespace),
	})
}

func (m *Resources) SetServiceName(serviceName string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_SERVICE_NAME,
		Value: attribute.StringValue(serviceName),
	})
}
