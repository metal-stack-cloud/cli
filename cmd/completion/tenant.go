package completion

import (
	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) AdminTenantListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv1.TenantServiceListRequest{}
	resp, err := c.Client.Adminv1().Tenant().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Tenants {
		names = append(names, s.Login+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
