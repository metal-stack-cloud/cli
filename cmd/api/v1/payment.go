package v1

import (
	"fmt"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func newPaymentCmd(c *config.Config) *cobra.Command {
	paymentCmd := &cobra.Command{
		Use:   "payment",
		Short: "manage payment of the metalstack.cloud",
	}

	showDefaultPricesCmd := &cobra.Command{
		Use:   "show-default-prices",
		Short: "show default prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			req := &apiv1.PaymentServiceGetDefaultPricesRequest{}

			resp, err := c.Client.Apiv1().Payment().GetDefaultPrices(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to list methods: %w", err)
			}

			return c.ListPrinter.Print(resp.Msg)
		},
	}

	paymentCmd.AddCommand(showDefaultPricesCmd)

	return paymentCmd
}
