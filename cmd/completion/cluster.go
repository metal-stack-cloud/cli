package completion

import (
	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
)

func (c *Completion) ClusterListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.ClusterServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv1().Cluster().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, c := range resp.Msg.Clusters {
		c := c
		names = append(names, c.Uuid+"\t"+c.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ClusterWorkerGroupsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.ClusterServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv1().Cluster().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, c := range resp.Msg.Clusters {
		c := c
		for _, w := range c.Workers {
			w := w
			names = append(names, w.Name)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ClusterPurposeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"production", "infrastructure", "evaluation"}, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ClusterAdminOperationCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"reconcile", "retry", "maintain"}, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AdminClusterListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func (c *Completion) AdminClusterNameListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv1.ClusterServiceListRequest{
		Project: pointer.PointerOrNil(c.Project),
	}
	resp, err := c.Client.Adminv1().Cluster().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Clusters {
		names = append(names, s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) AdminClusteFirewallListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	clusterID, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	req := &adminv1.ClusterServiceGetRequest{
		Uuid:         clusterID,
		WithMachines: true,
	}

	resp, err := c.Client.Adminv1().Cluster().Get(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, machine := range resp.Msg.Machines {
		machine := machine

		if machine.Role != "firewall" {
			continue
		}

		names = append(names, machine.Uuid+"\t"+machine.Hostname)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
