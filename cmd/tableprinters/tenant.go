package tableprinters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"

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

	for _, tenant := range data {
		id := tenant.Login
		name := tenant.Name
		email := tenant.Email
		admitted := strconv.FormatBool(tenant.Admitted)
		since := humanize.Time(tenant.CreatedAt.AsTime())
		provider := tenant.OauthProvider.Enum().String()
		coupons := "-"
		couponsWide := coupons
		if tenant.PaymentDetails != nil {
			cs := []string{}
			csw := []string{}
			for _, c := range tenant.PaymentDetails.Coupons {
				cs = append(cs, c.Name)
				csw = append(csw, fmt.Sprintf("%s %s", c.Name, c.CreatedAt.AsTime()))
			}
			coupons = strings.Join(cs, "\n")
			couponsWide = strings.Join(csw, "\n")
		}

		if wide {
			rows = append(rows, []string{id, name, email, provider, since, admitted, couponsWide})
		} else {
			rows = append(rows, []string{id, name, email, provider, since, admitted, coupons})
		}
	}

	return header, rows, nil
}

func (t *TablePrinter) TenantMemberTable(data []*apiv1.TenantMember, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Role", "Since"}

	for _, member := range data {
		row := []string{
			member.Id,
			member.Role.String(),
			humanize.Time(member.CreatedAt.AsTime()),
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}

func (t *TablePrinter) TenantInviteTable(data []*apiv1.TenantInvite, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"Secret", "Tenant", "Role", "Expires in"}

	for _, invite := range data {
		row := []string{
			invite.Secret,
			invite.Tenant,
			invite.Role.String(),
			humanize.Time(invite.ExpiresAt.AsTime()),
		}

		rows = append(rows, row)
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
