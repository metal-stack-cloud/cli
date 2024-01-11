package v1

import (
	"connectrpc.com/connect"
	v1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func newAssetCmd(c *config.Config) *cobra.Command {
	assetCmd := &cobra.Command{
		Use:   "asset",
		Short: "show asset",
		Long:  "assets are boundaries of consumable objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := c.Client.Apiv1().Asset().List(c.NewRequestContext(), connect.NewRequest(&v1.AssetServiceListRequest{}))
			if err != nil {
				return err
			}

			return c.ListPrinter.Print(resp.Msg.Assets)
		},
	}
	return assetCmd
}
