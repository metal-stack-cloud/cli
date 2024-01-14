package tableprinters

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) ClusterTable(data []*apiv1.Cluster, wide bool) ([]string, [][]string, error) {

	var (
		rows [][]string
		// header = []string{"ClusterStatus", "ID", "Name", "Project", "Kubernetes Version", "Nodes", "Uptime"}
		header = []string{"UID", "Tenant", "Project", "Name", "Version", "Partition", "Operation", "Progress", "Api", "Control", "Nodes", "Sys", "Size", "Age", "Purpose"}

		statusIcon = func(s string) string {
			switch s {
			case "True":
				return color.GreenString("✔")
			case "False":
				return color.RedString("✗")
			default:
				return color.YellowString("?")
			}
		}
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
			progress = fmt.Sprintf("%d%% [%s]", cluster.Status.Progress, cluster.Status.Type)
			api = statusIcon(cluster.Status.ApiServerReady)
			control = statusIcon(cluster.Status.ControlPlaneReady)
			nodes = statusIcon(cluster.Status.NodesReady)
			system = statusIcon(cluster.Status.SystemComponentsReady)
		}
		purpose := pointer.SafeDeref(cluster.Purpose)

		rows = append(rows, []string{
			cluster.Uuid,
			cluster.Tenant,
			cluster.Project,
			cluster.Name,
			cluster.Kubernetes.Version,
			cluster.Partition,
			operation,
			progress,
			api,
			control,
			nodes,
			system,
			nodesRange,
			humanize.Time(cluster.CreatedAt.AsTime()),
			purpose,
		})
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetColumnAlignment([]int{
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_CENTER,
			tablewriter.ALIGN_CENTER,
			tablewriter.ALIGN_CENTER,
			tablewriter.ALIGN_CENTER,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_LEFT,
		})
	})

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
