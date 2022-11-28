package v1

import (
	"github.com/bufbuild/connect-go"
	v1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func NewTokenCmd(c *config.Config) *cobra.Command {
	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "token commands",
		Long:  "token can be used to talk to our api",
		RunE: func(cmd *cobra.Command, args []string) error {
			tcr := &v1.TokenServiceCreateRequest{
				Subject: "get-pi",
				Permissions: []*v1.ProjectPermission{
					{
						Project: "p1",
						Permissions: []string{
							"Get",
						},
					},
				},
			}
			resp, err := c.Apiv1Client.Token().Create(c.Ctx, connect.NewRequest(tcr))
			if err != nil {
				return err
			}

			if err := c.Pf.NewPrinterDefaultYAML(c.Out).Print("Token:" + resp.Msg.Token); err != nil {
				return err
			}

			return nil
		},
	}
	return tokenCmd
}
