package main

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/go-faker/faker/v4"
	"net/http"
	"testProject/api"
)

// Register various metrics.
// Metric name may contain labels in Prometheus format - see below.
var (
	// Register counter without labels.
	requestsTotal = metrics.NewCounter("requests_total")

	// Register summary with a single label.
	requestDuration = metrics.NewSummary(`requests_duration_seconds{path="/foobar/baz"}`)

	// Register gauge with two labels.
	queueSize = metrics.NewGauge(fmt.Sprintf(`queue_size{cluster_id=%s,cluster_name=%s,net_type=%s,nodepool_id=%s,nodepool_name=%s,subnet_id=%s}`, "fake-id", "fake-string",
		"fake-string", "fake-string", "fake-string", "fake-string"), func() float64 {
		return float64(1)
	})

	// Register histogram with a single label.
	responseSize = metrics.NewHistogram(`response_size{path="/foo/bar"}`)
)

func convertStructToLabel(relation api.SubnetRelation) {
	// 节点池
	for _, subnet := range relation.NodePoolSubnets {
		for _, s := range subnet.SubnetID {
			metrics.NewGauge(fmt.Sprintf(`queue_size{cluster_id=%s,cluster_name=%s,net_type=%s,nodepool_id=%s,nodepool_name=%s,subnet_id=%s}`, relation.ClusterID, relation.ClusterName,
				"node", subnet.NodePoolID, subnet.NodePoolName, s), func() float64 {
				return float64(1)
			})
		}
	}
	// 容器网络

	// 自定义节点池
}
func main() {
	s := api.SubnetRelation{}
	err := faker.FakeData(&s)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	convertStructToLabel(s)
	metrics.GetDefaultSet()
	// Expose the registered metrics at `/metrics` path.
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	_ = http.ListenAndServe(":8080", nil)

}
