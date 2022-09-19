package printer

import (
	"sort"
	"strconv"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) CustomerTable(data []*apiv1.PaymentCustomer, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Admitted"}
	if wide {
		header = []string{"ID", "Name", "Email", "Admitted"}
	}

	sort.SliceStable(data, func(i, j int) bool { return *data[i].CustomerId < *data[j].CustomerId })
	for _, customer := range data {
		id := pointer.SafeDeref(customer.CustomerId)
		name := pointer.SafeDeref(customer.Name)
		email := pointer.SafeDeref(customer.Email)

		if wide {
			rows = append(rows, []string{id, name, email, strconv.FormatBool(customer.Admitted)})
		} else {
			rows = append(rows, []string{id, name, email, strconv.FormatBool(customer.Admitted)})
		}
	}

	return header, rows, nil
}
