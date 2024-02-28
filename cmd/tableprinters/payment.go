package tableprinters

import (
	"fmt"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/olekukonko/tablewriter"
)

func (t *TablePrinter) PaymentPricesTable(data []*apiv1.Price, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"Type", "Name", "Description", "Price", "Unit", "Currency"}

	for _, price := range data {
		description := price.GetDescription()
		if !wide {
			description = genericcli.TruncateEnd(description, 80)
		}

		row := []string{
			price.GetProductType().String(),
			price.GetName(),
			description,
			fmt.Sprintf("%f", price.GetUnitAmountDecimal()),
			price.GetUnitLabel(),
			price.GetCurrency(),
		}

		rows = append(rows, row)
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
