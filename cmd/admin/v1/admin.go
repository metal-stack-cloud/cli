package v1

import (
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func NewAdminCmd(c *config.Config) *cobra.Command {

	adminCmd := &cobra.Command{
		Use:          "admin",
		Short:        "admin commands",
		Long:         "",
		SilenceUsage: true,
		Hidden:       true,
	}

	adminCmd.AddCommand(newTenantCmd(c))
	adminCmd.AddCommand(newCouponCmd(c))
	adminCmd.AddCommand(newStorageCmd(c))
	adminCmd.AddCommand(newClusterCmd(c))

	return adminCmd
}
