package v1

import (
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
	cmd.AddCommand(newAssetCmd(c))
	cmd.AddCommand(newAuditCmd(c))
	cmd.AddCommand(newClusterCmd(c))
	cmd.AddCommand(newHealthCmd(c))
	cmd.AddCommand(newIPCmd(c))
	cmd.AddCommand(newMethodsCmd(c))
	cmd.AddCommand(newPaymentCmd(c))
	cmd.AddCommand(newProjectCmd(c))
	cmd.AddCommand(newStorageCmd(c))
	cmd.AddCommand(newTenantCmd(c))
	cmd.AddCommand(newTokenCmd(c))
	cmd.AddCommand(newUserCmd(c))
	cmd.AddCommand(newVersionCmd(c))
}

func newStorageCmd(c *config.Config) *cobra.Command {
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "storage commands",
		Long:  "volume and snapshot actions",
	}

	storageCmd.AddCommand(newVolumeCmd(c))
	storageCmd.AddCommand(newSnapshotCmd(c))

	return storageCmd
}
