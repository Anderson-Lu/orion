package tracing

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	_KEY_RESOURCE_ENV          = "deployment.environment"
	_KEY_RESOURCE_SERVICE_NAME = "service.name"
	_KEY_RESOURCE_NAMESPACE    = "service.namespace"
	_KEY_RESOURCE_INSTANCE_ID  = "service.instance.id"
	_KEY_RESOURCE_INSTANCE_IP  = "service.instance.ip"
)

type Resources struct {
	kvs []attribute.KeyValue
}

func (m *Resources) Env(env string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   _KEY_RESOURCE_ENV,
		Value: attribute.StringValue(env),
	})
}

func (m *Resources) InstanceId(id string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   _KEY_RESOURCE_INSTANCE_ID,
		Value: attribute.StringValue(id),
	})
}

func (m *Resources) IP(ipAddr string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   _KEY_RESOURCE_INSTANCE_IP,
		Value: attribute.StringValue(ipAddr),
	})
}

func (m *Resources) ServiceName(name string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   _KEY_RESOURCE_SERVICE_NAME,
		Value: attribute.StringValue(name),
	})
}

func (m *Resources) Namespace(namespace string) {
	m.kvs = append(m.kvs, attribute.KeyValue{
		Key:   _KEY_RESOURCE_NAMESPACE,
		Value: attribute.StringValue(namespace),
	})
}
