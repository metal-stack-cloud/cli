package admin

import (
	"fmt"

	v1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func NewCustomerCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customer",
		Short: "customer admin",
		Long:  "administrative commands for customers",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := c.Adminv1Client.Customer().List(c.Ctx, &v1.CustomerServiceListRequest{})
			if err != nil {
				return fmt.Errorf("failed to get health: %w", err)
			}

			return c.Pf.NewPrinter(c.Out).Print(resp.Customers)
		},
	}
	return cmd
}
