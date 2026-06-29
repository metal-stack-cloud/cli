package tableprinters

import (
	"fmt"
	"strconv"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
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

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) PaymentSubscriptionUsageTable(items []*apiv1.SubscriptionUsageItem, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Name", "Usage", "Start", "End"}

	for _, item := range items {
		row := []string{
			item.SubscriptionItemId,
			item.SubscriptionItemName,
			strconv.FormatInt(item.TotalUsage, 10),
			item.PeriodStart.String(),
			item.PeriodEnd.String(),
		}

		rows = append(rows, row)
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
