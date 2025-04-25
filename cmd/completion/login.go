package completion

import (
	"github.com/spf13/cobra"
)

func (c *Completion) LoginProviderCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"github", "azure", "google"}, cobra.ShellCompDirectiveNoFileComp
}
