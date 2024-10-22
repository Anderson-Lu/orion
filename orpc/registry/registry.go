package registry

import "context"

type RegisteyMode int

const (
	RegisteyConsul      RegisteyMode = 1 // default mode using consul
	RegisteyServiceMesh RegisteyMode = 5 // sidecar mode using service mesh
)

type Node struct {
	IP       string
	Port     uint32
	Weight   uint32
	UpdateAt int64
}

type IRegistryAddOption interface{}

type RegistryAddOptionTags []string

type IRegistryGetOption interface{}

type IRegistry interface {
	AddNode(ctx context.Context, ip string, port uint32, opts ...IRegistryAddOption) error
	GetNodes(ctx context.Context, name string) ([]*Node, error)
	RemoveNode(ctx context.Context) error
}
