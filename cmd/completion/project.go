package completion

import (
	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (c *Completion) ProjectListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.ProjectServiceListRequest{}
	resp, err := c.Client.Apiv1().Project().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.Msg.GetProjects() {
		names = append(names, s.Uuid+"\t"+s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ProjectRoleCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	for value, name := range apiv1.ProjectRole_name {
		if value == 0 {
			continue
		}

		names = append(names, name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ProjectInviteListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv1().Project().InvitesList(c.Ctx, connect.NewRequest(&apiv1.ProjectServiceInvitesListRequest{
		Project: c.Project,
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

func (c *Completion) ProjectMemberListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := c.Client.Apiv1().Project().Get(c.Ctx, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
		Project: c.Project,
	}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string

	for _, member := range resp.Msg.ProjectMembers {
		names = append(names, member.Id+"\t"+member.Role.String())
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
