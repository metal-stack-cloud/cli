package tableprinters

import (
	"fmt"
	"strconv"
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/pkg/helpers"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) TokenTable(data []*apiv1.Token, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"Type", "ID", "Admin", "User", "Description", "Roles", "Perms", "Expires"}

	for _, token := range data {
		expires := token.Expires.AsTime().Format(time.DateTime + " MST")
		expiresIn := helpers.HumanizeDuration(time.Until(token.Expires.AsTime()))
		admin := ""
		if token.AdminRole != nil {
			admin = token.AdminRole.String()
		}

		row := []string{
			token.TokenType.String(),
			token.Uuid,
			admin,
			token.UserId,
			token.Description,
			strconv.Itoa(len(token.TenantRoles) + len(token.ProjectRoles)),
			strconv.Itoa(len(token.Permissions)),
			fmt.Sprintf("%s (in %s)", expires, expiresIn),
		}

		rows = append(rows, row)
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
