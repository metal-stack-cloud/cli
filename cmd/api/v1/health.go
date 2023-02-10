package v1

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	v1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func NewHealthCmd(c *config.Config) *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "print the client and server health information",
		Long:  "print the client and server health information",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := c.Apiv1Client.Health().Get(c.Ctx, connect.NewRequest(&v1.HealthServiceGetRequest{}))
			if err != nil {
				return fmt.Errorf("failed to get health: %w", err)
			}

			return c.DescribePrinter.Print(resp.Msg.Health)
		},
	}
	return healthCmd
}
