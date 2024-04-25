package tableprinters

import (
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
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

func (t *TablePrinter) ProjectInviteTable(data []*apiv1.ProjectInvite, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"Secret", "Project", "Role", "Expires in"}

	for _, invite := range data {
		row := []string{
			invite.Secret,
			invite.Project,
			invite.Role.String(),
			humanize.Time(invite.ExpiresAt.AsTime()),
		}

		rows = append(rows, row)
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}

func (t *TablePrinter) ProjectMemberTable(data []*apiv1.ProjectMember, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Role", "Inherited", "Since"}

	for _, member := range data {
		row := []string{
			member.Id,
			member.Role.String(),
			strconv.FormatBool(member.ImplicitMembership),
			humanize.Time(member.CreatedAt.AsTime()),
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}
