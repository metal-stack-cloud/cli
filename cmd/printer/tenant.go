package printer

import (
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) TenantTable(data []*apiv1.Tenant, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Registered", "Admitted"}
	if wide {
		header = []string{"ID", "Name", "Email", "Registered", "Admitted"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Login < data[j].Login })
	for _, tenant := range data {
		id := tenant.Login
		name := tenant.Name
		email := tenant.Email
		admitted := strconv.FormatBool(tenant.Admitted)
		since := humanize.Time(tenant.CreatedAt.AsTime())

		if wide {
			rows = append(rows, []string{id, name, email, since, admitted})
		} else {
			rows = append(rows, []string{id, name, email, since, admitted})
		}
	}

	return header, rows, nil
}
