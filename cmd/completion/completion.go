package completion

import (
	"context"

	"github.com/metal-stack-cloud/api/go/client"
	"github.com/spf13/cobra"
)

type Completion struct {
	Apiv1Client   client.Apiv1
	Adminv1Client client.Adminv1
	Ctx           context.Context
}

func OutputFormatListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp
}
