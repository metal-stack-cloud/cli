package completion

import (
	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) TenantListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.TenantServiceListRequest{}
	resp, err := c.Client.Apiv1().Tenant().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, t := range resp.Msg.Tenants {
		names = append(names, t.Login+"\t"+t.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantRoleCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for value, name := range apiv1.TenantRole_name {
		if value == 0 {
			continue
		}

		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantOauthProviderCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var providers []string

	for _, name := range apiv1.OAuthProvider_name {
		providers = append(providers, name)
	}

	return providers, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantInviteListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectResp, err := c.Client.Apiv1().Project().Get(c.Ctx, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
		Project: c.Project,
	}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	resp, err := c.Client.Apiv1().Tenant().InvitesList(c.Ctx, connect.NewRequest(&apiv1.TenantServiceInvitesListRequest{
		Login: projectResp.Msg.Project.Tenant,
	}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, invite := range resp.Msg.Invites {
		names = append(names, invite.Secret+"\t"+invite.Role.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TenantMemberListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectResp, err := c.Client.Apiv1().Project().Get(c.Ctx, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
		Project: c.Project,
	}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	resp, err := c.Client.Apiv1().Tenant().Get(c.Ctx, connect.NewRequest(&apiv1.TenantServiceGetRequest{
		Login: projectResp.Msg.Project.Tenant,
	}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, member := range resp.Msg.TenantMembers {
		names = append(names, member.Id+"\t"+member.Role.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

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
