package v1

import (
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
	cmd.AddCommand(newVersionCmd(c))
	cmd.AddCommand(newHealthCmd(c))
	cmd.AddCommand(newAssetCmd(c))
	cmd.AddCommand(newTokenCmd(c))
	cmd.AddCommand(newIPCmd(c))
	cmd.AddCommand(newStorageCmd(c))
	cmd.AddCommand(newClusterCmd(c))
	cmd.AddCommand(newProjectCmd(c))
}
