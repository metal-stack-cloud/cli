package completion

import (
	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) VolumeListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.VolumeServiceListRequest{
		Project: c.Project,
	}
	resp, err := c.Client.Apiv1().Volume().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, c := range resp.Msg.Volumes {
		c := c
		names = append(names, c.Uuid+"\t"+c.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
