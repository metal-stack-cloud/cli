package completion

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) ClusterListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv1.ClusterServiceListRequest{}
	resp, err := c.Client.Adminv1().Cluster().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Clusters {
		fmt.Println(s.Uuid)
		names = append(names, s.Uuid)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) ClusterPurposeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"production", "infrastructure", "evaluation"}, cobra.ShellCompDirectiveNoFileComp
}
