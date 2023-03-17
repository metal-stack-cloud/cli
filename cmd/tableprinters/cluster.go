package tableprinters

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) ClusterTable(data []*apiv1.Cluster, wide bool) ([]string, [][]string, error) {

	var (
		rows   [][]string
		header = []string{"ClusterStatus", "ID", "Name", "Project", "Kubernetes Version", "Nodes", "Uptime"}
	)

	for _, cluster := range data {
		var clusterState string
		switch cluster.Status.State {
		case "Error":
			clusterState = color.RedString(dot)
		case "Processing":
			processing := color.New(color.FgHiRed, color.FgHiYellow).SprintFunc()
			clusterState = processing(dot)
		case "Success":
			clusterState = color.GreenString(dot)
		}

		var totalMinNodes, totalMaxNodes uint32

		for _, worker := range cluster.Workers {
			totalMinNodes += worker.Minsize
			totalMaxNodes += worker.Maxsize
		}
		nodesRange := fmt.Sprintf("%v - %v", totalMinNodes, totalMaxNodes)

		rows = append(rows, []string{clusterState, cluster.Uuid, cluster.Name, cluster.Project, cluster.Kubernetes.Version, nodesRange, humanize.Time(cluster.CreatedAt.AsTime())})
	}

	return header, rows, nil

}
