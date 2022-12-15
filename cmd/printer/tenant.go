package printer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) TenantTable(data []*apiv1.Tenant, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Provider", "Registered", "Admitted", "Coupons"}
	if wide {
		header = []string{"ID", "Name", "Email", "Provider", "Registered", "Admitted", "Coupons"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Login < data[j].Login })
	for _, tenant := range data {
		id := tenant.Login
		name := tenant.Name
		email := tenant.Email
		admitted := strconv.FormatBool(tenant.Admitted)
		since := humanize.Time(tenant.CreatedAt.AsTime())
		provider := tenant.OauthProvider.Enum().String()
		coupons := "-"
		if tenant.PaymentDetails != nil {
			cs := []string{}
			for _, c := range tenant.PaymentDetails.Coupons {
				cs = append(cs, fmt.Sprintf("%.2f%s til %s", float64(c.AmountOff/100), c.Currency, c.RedeemBy))
			}
			coupons = strings.Join(cs, "\n")
		}

		if wide {
			rows = append(rows, []string{id, name, email, provider, since, admitted, coupons})
		} else {
			rows = append(rows, []string{id, name, email, provider, since, admitted, coupons})
		}
	}

	return header, rows, nil
}
