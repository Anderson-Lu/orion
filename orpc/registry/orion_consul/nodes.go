package orion_consul

type OrionNodesNotifyFunc func(nodes []OrionNode)

type OrionNode struct {
	Namespace   string
	Datacenter  string
	Host        string
	Port        int
	ServiceName string
	Tags        []string
}
