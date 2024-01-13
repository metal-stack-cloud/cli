package completion

import (
	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
)

func (c *Completion) ClusterAdminListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv1.ClusterServiceListRequest{
		Project: pointer.PointerOrNil(c.Project),
	}
	resp, err := c.Client.Adminv1().Cluster().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Clusters {
		names = append(names, s.Uuid+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ClusterListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.ClusterServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv1().Cluster().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Clusters {
		names = append(names, s.Uuid+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ClusterPurposeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"production", "infrastructure", "evaluation"}, cobra.ShellCompDirectiveNoFileComp
}
