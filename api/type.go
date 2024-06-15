package api

type SubnetRelation struct {
	ClusterID        string
	ClusterName      string
	NetType          string
	NodePoolSubnets  []NodePoolSubnet
	ContainerSubnets []ContainerSubnet
}

type NodePoolSubnet struct {
	NodePoolName string
	NodePoolID   string
	SubnetID     []string
}

type ContainerSubnet struct {
	// 只能是容器子网和自定义子网中的一项
	IsCustomContainerSubnet bool
	ContainerSubnetID       string
	CustomContainerSubnet   CustomContainerSubnet
}
type CustomContainerSubnet struct {
	CustomContainerSubnetName string
	SubnetID                  []string
}

type MyMetric struct {
}
