package printer

import (
	"fmt"
	"sort"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) CouponTable(data []*apiv1.Coupon, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "AmountOff", "Duration"}
	if wide {
		header = []string{"ID", "Name", "AmountOff", "Duration"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Id < data[j].Id })
	for _, coupon := range data {
		id := coupon.Id
		name := coupon.Name
		amount := fmt.Sprintf("%.2f %s", float64(coupon.AmountOff/100), coupon.Currency)
		duration := fmt.Sprintf("%d month", coupon.DurationInMonth)

		if wide {
			rows = append(rows, []string{id, name, amount, duration})
		} else {
			rows = append(rows, []string{id, name, amount, duration})
		}
	}

	return header, rows, nil
}
