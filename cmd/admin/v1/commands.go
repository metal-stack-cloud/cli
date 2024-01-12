package v1

import (
	"fmt"

	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddCmds(cmd *cobra.Command, c *config.Config) {
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
	adminCmd.AddCommand(newTokenCmd(c))

	cmd.AddCommand(adminCmd)
}

func newStorageCmd(c *config.Config) *cobra.Command {
	storageCmd := &cobra.Command{
		Use:          "storage",
		Short:        "storage commands",
		Long:         "",
		SilenceUsage: true,
		Hidden:       true,
	}

	clusterInfoCmd := &cobra.Command{
		Use:   "clusterinfo",
		Short: "storage clusterinfo",
		Long:  "show detailed info about the storage cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			req := &adminv1.StorageServiceClusterInfoRequest{}
			if viper.IsSet("partition") {
				req.Partition = pointer.Pointer(viper.GetString("partition"))
			}

			resp, err := c.Client.Adminv1().Storage().ClusterInfo(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to get clusterinfo: %w", err)
			}

			return c.ListPrinter.Print(resp.Msg.Infos)
		},
	}
	clusterInfoCmd.Flags().StringP("partition", "", "", "optional partition to filter for storage cluster info")

	storageCmd.AddCommand(newVolumeCmd(c))
	storageCmd.AddCommand(newSnapshotCmd(c))
	storageCmd.AddCommand(clusterInfoCmd)

	return storageCmd
}
