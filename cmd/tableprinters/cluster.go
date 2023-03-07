package tableprinters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) ClusterTable(data []*apiv1.Cluster, wide bool) ([]string, [][]string, error) {

	var (
		rows   [][]string
		header = []string{"ClusterStatus", "ID", "Name", "Project", "KubernetesSpec"} //TODO: Which fields to display?
	)

	for _, cluster := range data {
		// var s string
		// switch cluster.Status {
		// 	case apiv1.ClusterStatus.
		// }
		rows = append(rows, []string{cluster.Status.State, cluster.Uuid, cluster.Name, cluster.Project, cluster.Kubernetes.Version})
	}

	return header, rows, nil

}
