package completion

import (
	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) IpListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.IPServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv1().IP().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, s := range resp.Msg.Ips {
		names = append(names, s.Uuid+"\t"+s.Ip+"\t"+s.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
