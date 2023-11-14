package v1

import (
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

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
