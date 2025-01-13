package completion

import (
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (c *Completion) TokenListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.TokenServiceListRequest{}
	resp, err := c.Client.Apiv1().Token().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.Msg.Tokens {
		fmt.Println(s.Uuid)
		names = append(names, s.Uuid+"\t"+s.Description)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenProjectRolesCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv1().Method().TokenScopedList(c.Ctx, connect.NewRequest(&apiv1.MethodServiceTokenScopedListRequest{}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var roles []string

	for project, role := range methods.Msg.ProjectRoles {
		roles = append(roles, project+"="+role.String())
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenTenantRolesCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv1().Method().TokenScopedList(c.Ctx, connect.NewRequest(&apiv1.MethodServiceTokenScopedListRequest{}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var roles []string

	for tenant, role := range methods.Msg.TenantRoles {
		roles = append(roles, tenant+"="+role.String())
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenAdminRoleCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var roles []string

	for _, role := range apiv1.AdminRole_name {
		roles = append(roles, role)
	}

	return roles, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) TokenPermissionsCompletionfunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	methods, err := c.Client.Apiv1().Method().TokenScopedList(c.Ctx, connect.NewRequest(&apiv1.MethodServiceTokenScopedListRequest{}))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	subject := ""
	if s, _, ok := strings.Cut(toComplete, "="); ok {
		subject = s
	}

	if subject == "" {
		var perms []string

		for _, p := range methods.Msg.Permissions {
			perms = append(perms, p.Subject)
		}

		return perms, cobra.ShellCompDirectiveNoFileComp
	}

	// FIXME: completion does not work at this point, investigate why

	var perms []string

	for _, p := range methods.Msg.Permissions {
		perms = append(perms, p.Methods...)
	}

	return perms, cobra.ShellCompDirectiveDefault
}
