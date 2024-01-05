package completion

import (
	"fmt"

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
		fmt.Println(s.Uuid)
		names = append(names, s.Uuid+"\t"+s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
