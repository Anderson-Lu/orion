package tracing

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	KEY_RESOURCE_ENV          = "deployment.environment"
	KEY_RESOURCE_NAMESPACE    = "deployment.namespace"
	KEY_RESOURCE_SERVICE_NAME = "service.name"
	KEY_RESOURCE_INSTANCE_ID  = "service.instance.id"
	KEY_RESOURCE_INSTANCE_IP  = "service.instance.ip"
	KEY_SPAN_ERRCODE          = "span.code"
	KEY_SPAN_KIND             = "span.kind"
)

type Resources struct {
	kvs []attribute.KeyValue
}

func (m *Resources) Env(env string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_ENV,
		Value: attribute.StringValue(env),
	})
}

func (m *Resources) InstanceId(id string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_INSTANCE_ID,
		Value: attribute.StringValue(id),
	})
}

func (m *Resources) IP(ipAddr string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_INSTANCE_IP,
		Value: attribute.StringValue(ipAddr),
	})
}

func (m *Resources) ServiceName(name string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_SERVICE_NAME,
		Value: attribute.StringValue(name),
	})
}

func (m *Resources) Namespace(namespace string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_RESOURCE_NAMESPACE,
		Value: attribute.StringValue(namespace),
	})
}

func (m *Resources) KindClient() {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   KEY_SPAN_KIND,
		Value: attribute.StringValue("client"),
	})
}
