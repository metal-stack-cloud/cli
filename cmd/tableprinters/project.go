package tableprinters

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) ProjectTable(data []*apiv1.Project, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Tenant", "Name", "Description", "Creation Date"}

	for _, project := range data {
		row := []string{
			project.Uuid,
			project.Tenant,
			project.Name,
			genericcli.TruncateEnd(project.Description, 80),
			project.CreatedAt.AsTime().Format(time.DateTime + " MST"),
		}

		rows = append(rows, row)
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
