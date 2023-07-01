package v1

import (
	"github.com/bufbuild/connect-go"
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
			resp, err := c.Client.Apiv1().Asset().List(c.Ctx, connect.NewRequest(&v1.AssetServiceListRequest{}))
			if err != nil {
				return err
			}

			if err := c.ListPrinter.Print(resp.Msg); err != nil {
				return err
			}

			return nil
		},
	}
	return assetCmd
}
