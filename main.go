package main

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/go-faker/faker/v4"
	"net/http"
	"testProject/api"
)

func convertStructToLabel(relation api.SubnetRelation) []api.MyMetric {
	var subnetMetrics []api.MyMetric
	// 节点池
	for _, nodePoolSubnet := range relation.NodePoolSubnets {
		for _, id := range nodePoolSubnet.SubnetID {
			m := api.MyMetric{}
			m.ClusterID = relation.ClusterID
			m.ClusterName = relation.ClusterName
			m.NetType = api.NodePoolSubnetType
			m.NodePoolID = nodePoolSubnet.NodePoolID
			m.NodePoolName = nodePoolSubnet.NodePoolName
			m.SubnetID = id
			subnetMetrics = append(subnetMetrics, m)
		}
	}
	// 容器网络和自定义容器网络
	for _, containerSubnet := range relation.ContainerSubnets {
		if containerSubnet.IsCustomContainerSubnet {
			m := api.MyMetric{}
			m.ClusterID = relation.ClusterID
			m.ClusterName = relation.ClusterName
			m.NetType = api.ContainerSubnetType
			m.NodePoolID = "-1"
			m.NodePoolName = "NotNodePool"
			m.SubnetID = containerSubnet.ContainerSubnetID
			subnetMetrics = append(subnetMetrics, m)
		} else {
			for _, customContainerSubnetID := range containerSubnet.CustomContainerSubnet.SubnetID {
				m := api.MyMetric{}
				m.ClusterID = relation.ClusterID
				m.ClusterName = relation.ClusterName
				m.NetType = api.CustomContainerSubnetType
				m.NodePoolID = "-1"
				m.NodePoolName = "NotNodePool"
				m.SubnetID = customContainerSubnetID
				subnetMetrics = append(subnetMetrics, m)
			}
		}
	}
	return subnetMetrics
}

func GenerateMetric(myMetrics []api.MyMetric) {
	for _, metric := range myMetrics {
		metrics.NewGauge(fmt.Sprintf(`subnet_relation{cluster_id="%s",cluster_name="%s",net_type="%s",nodepool_id="%s",nodepool_name="%s",subnet_id="%s"}`, metric.ClusterID, metric.ClusterName,
			"node", metric.NodePoolID, metric.NodePoolName, metric.SubnetID), func() float64 {
			return float64(1)
		})
	}

}
func main() {
	s := api.SubnetRelation{}
	err := faker.FakeData(&s)
	if err != nil {
		panic(err)
	}
	m := convertStructToLabel(s)
	GenerateMetric(m)
	// Expose the registered metrics at `/metrics` path.
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	_ = http.ListenAndServe(":8080", nil)

}
