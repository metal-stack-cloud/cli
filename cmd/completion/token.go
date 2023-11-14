package completion

import (
	"fmt"

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
