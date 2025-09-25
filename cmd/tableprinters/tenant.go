package tableprinters

import (
	"strconv"

	"github.com/dustin/go-humanize"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) TenantTable(data []*apiv1.Tenant, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Provider", "Registered", "Admitted", "Terms And Conditions"}
	if wide {
		header = []string{"ID", "Name", "Email", "Provider", "Registered", "Admitted", "Terms And Conditions"}
	}

	for _, tenant := range data {
		var (
			id                 = tenant.Login
			name               = tenant.Name
			email              = tenant.Email
			admitted           = strconv.FormatBool(tenant.Admitted)
			since              = humanize.Time(tenant.CreatedAt.AsTime())
			provider           = tenant.OauthProvider.Enum().String()
			termsAndConditions string
		)

		if tenant.TermsAndConditions != nil {
			termsAndConditions = strconv.FormatBool(tenant.TermsAndConditions.Accepted)
		}

		if wide {
			rows = append(rows, []string{id, name, email, provider, since, admitted, termsAndConditions})
		} else {
			rows = append(rows, []string{id, name, email, provider, since, admitted, termsAndConditions})
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
	header := []string{"Secret", "Tenant", "Invited By", "Role", "Expires in"}

	for _, invite := range data {
		row := []string{
			invite.Secret,
			invite.TargetTenant,
			invite.Tenant,
			invite.Role.String(),
			humanize.Time(invite.ExpiresAt.AsTime()),
		}

		rows = append(rows, row)
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
