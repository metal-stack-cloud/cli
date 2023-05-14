package tableprinters

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) ClusterTable(data []*apiv1.Cluster, wide bool) ([]string, [][]string, error) {

	var (
		rows [][]string
		// header = []string{"ClusterStatus", "ID", "Name", "Project", "Kubernetes Version", "Nodes", "Uptime"}
		header = []string{"UID", "Tenant", "Project", "Name", "Version", "Partition", "Operation", "Progress", "Api", "Control", "Nodes", "System", "Size", "Age", "Purpose"}
	)

	for _, cluster := range data {
		var totalMinNodes, totalMaxNodes uint32

		for _, worker := range cluster.Workers {
			totalMinNodes += worker.Minsize
			totalMaxNodes += worker.Maxsize
		}
		nodesRange := fmt.Sprintf("%v - %v", totalMinNodes, totalMaxNodes)

		api := ""
		control := ""
		nodes := ""
		system := ""
		operation := ""
		progress := "0%"
		if cluster.Status != nil {
			operation = cluster.Status.State
			progress = fmt.Sprintf("%d%% [%s]", cluster.Status.Progress, cluster.Status.State)
			api = cluster.Status.ApiServerReady
			control = cluster.Status.ControlPlaneReady
			nodes = cluster.Status.NodesReady
			system = cluster.Status.SystemComponentsReady
		}
		purpose := pointer.SafeDeref(cluster.Purpose)

		rows = append(rows, []string{
			cluster.Uuid, cluster.Tenant, cluster.Project, cluster.Name, cluster.Kubernetes.Version, cluster.Kubernetes.Version, operation, progress, api, control, nodes, system, nodesRange, humanize.Time(cluster.CreatedAt.AsTime()), purpose})
	}

	return header, rows, nil

}

func (t *TablePrinter) ClusterStatusLastErrorTable(data []*apiv1.ClusterStatusLastError, wide bool) ([]string, [][]string, error) {

	var (
		rows   [][]string
		header = []string{"Time", "Description", "Codes", "Task"}
	)

	for _, e := range data {
		rows = append(rows, []string{
			e.LastUpdateTime.String(),
			e.Description,
			strings.Join(e.Codes, ","),
			*e.TaskId,
		},
		)
	}

	return header, rows, nil

}
