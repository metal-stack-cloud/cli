package tableprinters

import (
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) CouponTable(data []*apiv1.Coupon, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "AmountOff", "Duration", "Redeemed", "Created"}
	if wide {
		header = []string{"ID", "Name", "AmountOff", "Duration", "Redeemed", "Created"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Id < data[j].Id })
	for _, coupon := range data {
		id := coupon.Id
		name := coupon.Name
		amount := fmt.Sprintf("%.2f %s", float64(coupon.AmountOff/100), coupon.Currency)
		duration := fmt.Sprintf("%d month", coupon.DurationInMonth)
		redeemed := fmt.Sprintf("%d/%d", coupon.TimesRedeemed, coupon.MaxRedemptions)
		created := humanize.Time(coupon.CreatedAt.AsTime())

		if wide {
			rows = append(rows, []string{id, name, amount, duration, redeemed, created})
		} else {
			rows = append(rows, []string{id, name, amount, duration, redeemed, created})
		}
	}

	return header, rows, nil
}
