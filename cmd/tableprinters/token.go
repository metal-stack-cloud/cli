package tableprinters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/pkg/helpers"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) TokenTable(data []*apiv1.Token, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Admin", "User", "Description", "Roles", "Perms", "Expires"}

	for _, token := range data {
		expires := token.Expires.AsTime().Format(time.DateTime + " MST")
		expiresIn := helpers.HumanizeDuration(time.Until(token.Expires.AsTime()))
		admin := isAdminToken(token)

		row := []string{
			token.Uuid,
			strconv.FormatBool(admin),
			token.UserId,
			token.Description,
			strconv.Itoa(len(token.Roles)),
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

func isAdminToken(token *apiv1.Token) bool {
	// TODO: maybe it would make more sense to put this information into the token itself?

	for _, role := range token.Roles {
		if role.Subject == "*" {
			return true
		}
	}

	// if there is any admin method contained in the permissions, we assume the token comes from an admin
	for _, perms := range token.Permissions {
		for _, perm := range perms.Methods {
			if strings.HasPrefix(perm, "/admin.v1") {
				return true
			}
		}
	}

	return false
}
