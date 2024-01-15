package tableprinters

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) ClusterTable(clusters []*apiv1.Cluster, machines map[string][]*adminv1.Machine, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"", "Tenant", "Project", "ID", "Name", "Partition", "Version", "Size", "Age"}

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
		statusShort = func(status *apiv1.ClusterStatus) string {
			progress := fmt.Sprintf("%d%%", status.Progress)

			if status.ApiServerReady == "True" && status.ControlPlaneReady == "True" && status.NodesReady == "True" && status.SystemComponentsReady == "True" {
				progress = color.GreenString(progress)
			} else {
				progress = color.YellowString(progress)
			}

			return progress
		}
	)

	if wide {
		header = []string{"ID", "Tenant", "Project", "Name", "Partition", "Purpose", "Version", "Operation", "Progress", "Api", "Control", "Nodes", "Sys", "Size", "Age"}
	}

	for _, cluster := range clusters {
		cluster := cluster

		var (
			totalMinNodes, totalMaxNodes uint32
		)

		for _, worker := range cluster.Workers {
			totalMinNodes += worker.Minsize
			totalMaxNodes += worker.Maxsize
		}

		nodesRange := fmt.Sprintf("%v - %v", totalMinNodes, totalMaxNodes)

		if wide {
			var (
				purpose   = pointer.SafeDeref(cluster.Purpose)
				api       = ""
				control   = ""
				nodes     = ""
				system    = ""
				operation = ""
				progress  = "0%"
			)

			if cluster.Status != nil {
				operation = cluster.Status.State
				progress = fmt.Sprintf("%d%% [%s]", cluster.Status.Progress, cluster.Status.Type)
				api = statusIcon(cluster.Status.ApiServerReady)
				control = statusIcon(cluster.Status.ControlPlaneReady)
				nodes = statusIcon(cluster.Status.NodesReady)
				system = statusIcon(cluster.Status.SystemComponentsReady)
			}

			rows = append(rows, []string{
				cluster.Uuid,
				cluster.Tenant,
				cluster.Project,
				cluster.Name,
				cluster.Partition,
				purpose,
				cluster.Kubernetes.Version,
				operation,
				progress,
				api,
				control,
				nodes,
				system,
				nodesRange,
				humanize.Time(cluster.CreatedAt.AsTime()),
			})
		} else {
			rows = append(rows, []string{
				statusShort(cluster.Status),
				cluster.Tenant,
				cluster.Project,
				cluster.Uuid,
				cluster.Name,
				cluster.Partition,
				cluster.Kubernetes.Version,
				nodesRange,
				humanize.Time(cluster.CreatedAt.AsTime()),
			})

			for i, machine := range machines[cluster.Uuid] {
				machine := machine

				prefix := "├"
				if i == len(machines[cluster.Uuid])-1 {
					prefix = "└"
				}
				prefix += "─╴"

				status := machine.Liveliness
				switch status {
				case "Alive":
					status = color.GreenString("✔")
				default:
					status = color.RedString("✗")
				}

				rows = append(rows, []string{
					status,
					"",
					"",
					prefix + machine.Uuid,
					machine.Hostname,
					machine.Partition,
					machine.Image,
					machine.Size,
					humanize.Time(machine.Created.AsTime()),
				})
			}
		}
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		if wide {
			table.SetColumnAlignment([]int{
				tablewriter.ALIGN_LEFT,
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
			})
		}
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

func (t *TablePrinter) ClusterMachineTable(data []*adminv1.ClusterServiceGetResponse, wide bool) ([]string, [][]string, error) {
	var (
		clusters []*apiv1.Cluster
		machines = map[string][]*adminv1.Machine{}
	)

	for _, cluster := range data {
		cluster := cluster

		clusters = append(clusters, cluster.Cluster)
		machines[cluster.Cluster.Uuid] = append(machines[cluster.Cluster.Uuid], cluster.Machines...)
	}

	return t.ClusterTable(clusters, machines, wide)
}
