package v1

import (
	"fmt"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type payment struct {
	c *config.Config
}

func newPaymentCmd(c *config.Config) *cobra.Command {
	w := &payment{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.PaymentCustomer]{
		BinaryName:  config.BinaryName,
		GenericCLI:  genericcli.NewGenericCLI(w).WithFS(c.Fs),
		Singular:    "payment",
		Plural:      "payments",
		Description: "manage payment of the metalstack.cloud",
		// Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
		},
		OnlyCmds: genericcli.OnlyCmds(
			genericcli.DescribeCmd,
			genericcli.CreateCmd,
			genericcli.UpdateCmd,
			genericcli.DeleteCmd,
			genericcli.ApplyCmd,
			genericcli.EditCmd,
		),
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant")
			cmd.Flags().String("email", "", "the email of the tenant")
			cmd.Flags().String("phone", "", "the phone number of the tenant")
			cmd.Flags().String("vat", "", "the vat of the tenant")

			cmd.Flags().String("stripe-public-token", config.DefaultStripePubToken, "the stripe public token")
		},
		CreateRequestFromCLI: func() (any, error) {
			// tenant, err := w.c.GetTenant()
			// if err != nil {
			// 	return nil, err
			// }

			// params := &stripe.PaymentMethodParams{
			// 	Type: stripe.String(string(stripe.PaymentMethodTypeCard)),
			// 	Card: &stripe.PaymentMethodCardParams{
			// 		Token: stripe.String("tok_de"),
			// 	},
			// }
			// result, err := paymentmethod.New(params)
			// if err != nil {
			// 	t.Errorf("Failed to create payment method: %v", err)
			// 	return
			// }

			// return &apiv1.PaymentServiceCreateRequest{
			// 	Login:           tenant,
			// 	Name:            viper.GetString("name"),
			// 	PaymentMethodId: "", // TODO
			// 	Email:           viper.GetString("email"),
			// 	Address:         &apiv1.Address{}, // TODO
			// 	Vat:             viper.GetString("vat"),
			// 	PhoneNumber:     pointer.PointerOrNil(viper.GetString("phone")),
			// }, nil
			return nil, nil
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant to update")
			cmd.Flags().String("description", "", "the description of the tenant to update")
		},
		// UpdateRequestFromCLI: w.updateRequestFromCLI,
		ValidArgsFn: w.c.Completion.TenantListCompletion,
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

			prices := resp.Msg.GetPrices()

			err = sorters.PriceSorter().SortBy(prices)
			if err != nil {
				return err
			}

			return c.ListPrinter.Print(prices)
		},
	}

	genericcli.AddSortFlag(showDefaultPricesCmd, sorters.PriceSorter())

	getSubscriptionUsageCmd := &cobra.Command{
		Use:   "get-subscription-usage",
		Short: "get-subscription-usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.GetSubscriptionUsage()
		},
	}

	getSubscriptionUsageCmd.Flags().String("tenant", "", "tenant of the subscription usage.")
	genericcli.Must(getSubscriptionUsageCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	return genericcli.NewCmds(cmdsConfig, showDefaultPricesCmd, getSubscriptionUsageCmd)
}

func (p *payment) Convert(r *apiv1.PaymentCustomer) (string, any, any, error) {
	panic("unimplemented")
}

func (p *payment) Create(rq any) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")

	// ctx, cancel := p.c.NewRequestContext()
	// defer cancel()

	// payment, err := p.c.Client.Apiv1().Payment().Create(ctx, connect.NewRequest(rq))
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to create payment data: %w", err)
	// }

	// return payment.Msg.GetCustomer(), nil
}

func (p *payment) Delete(id string) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")
}

func (p *payment) Get(id string) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")

	// ctx, cancel := p.c.NewRequestContext()
	// defer cancel()

	// tenant, err := p.c.GetTenant()
	// if err != nil {
	// 	return nil, err
	// }

	// payment, err := p.c.Client.Apiv1().Payment().Get(ctx, connect.NewRequest(&apiv1.PaymentServiceGetRequest{
	// 	Login: tenant,
	// }))
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to get payment data: %w", err)
	// }

	// return payment.Msg.GetCustomer(), nil
}

func (p *payment) List() ([]*apiv1.PaymentCustomer, error) {
	panic("unimplemented")
}

func (p *payment) Update(rq any) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")

	// ctx, cancel := p.c.NewRequestContext()
	// defer cancel()

	// payment, err := p.c.Client.Apiv1().Payment().Update(ctx, connect.NewRequest(rq))
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to update payment data: %w", err)
	// }

	// return payment.Msg.GetCustomer(), nil
}

func (p *payment) GetSubscriptionUsage() error {
	ctx, cancel := p.c.NewRequestContext()
	defer cancel()

	tenant, err := p.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := p.c.Client.Apiv1().Payment().GetSubscriptionUsage(ctx, connect.NewRequest(&apiv1.PaymentServiceGetSubscriptionUsageRequest{
		Login: tenant,
	}))
	if err != nil {
		return fmt.Errorf("unable to get subscription usage: %w", err)
	}

	return p.c.ListPrinter.Print(resp.Msg.GetSubscriptionUsageItems())
}
