package cmd

import (
	"fmt"

	v1 "github.com/metal-stack-cloud/api/go/v1"
	"github.com/spf13/cobra"
)

func newHealthCmd(c *config) *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "print the client and server health information",
		Long:  "print the client and server health information",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := c.client.Health().Get(c.ctx, &v1.HealthServiceGetRequest{})
			if err != nil {
				return fmt.Errorf("failed to get health: %w", err)
			}

			return c.pf.newPrinterDefaultYAML().Print(resp.Health)
		},
	}
	return healthCmd
}
