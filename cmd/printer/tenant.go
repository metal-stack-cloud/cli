package printer

import (
	"sort"
	"strconv"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) TenantTable(data []*apiv1.Tenant, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Admitted"}
	if wide {
		header = []string{"ID", "Name", "Email", "Admitted"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Login < data[j].Login })
	for _, user := range data {
		id := user.Login
		name := user.Name
		email := user.Email
		admitted := strconv.FormatBool(user.Admitted)

		if wide {
			rows = append(rows, []string{id, name, email, admitted})
		} else {
			rows = append(rows, []string{id, name, email, admitted})
		}
	}

	return header, rows, nil
}
