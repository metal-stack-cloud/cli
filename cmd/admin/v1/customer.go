package v1

import (
	"fmt"

	v1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type customer struct {
	c *config.Config
}

func NewCustomerCmd(c *config.Config) *cobra.Command {
	w := &customer{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.PaymentCustomer]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv1.PaymentCustomer](w).WithFS(c.Fs),
		Singular:        "customer",
		Plural:          "customers",
		Description:     "a customer of metal-stack cloud",
		Sorter:          sorters.CustomerSorter(),
		DescribePrinter: func() printers.Printer { return c.Pf.NewPrinterDefaultYAML(c.Out) },
		ListPrinter:     func() printers.Printer { return c.Pf.NewPrinter(c.Out) },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
	}

	admitCmd := &cobra.Command{
		Use:   "admit",
		Short: "admit a customer",
		Long:  "only admitted customers are allowed to consume resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			resp, err := c.Adminv1Client.Customer().Admit(c.Ctx, &v1.CustomerServiceAdmitRequest{
				CustomerId: id,
			})
			if err != nil {
				return fmt.Errorf("failed to admit customer: %w", err)
			}

			return c.Pf.NewPrinter(c.Out).Print(resp.Customer)
		},
	}

	return genericcli.NewCmds(cmdsConfig, admitCmd)
}

// Create implements genericcli.CRUD
func (c *customer) Create(rq any) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (c *customer) Delete(id string) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD
func (c *customer) Get(id string) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")
}

// List implements genericcli.CRUD
func (c *customer) List() ([]*apiv1.PaymentCustomer, error) {
	resp, err := c.c.Adminv1Client.Customer().List(c.c.Ctx, &v1.CustomerServiceListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return resp.Customers, nil
}

// ToCreate implements genericcli.CRUD
func (c *customer) ToCreate(r *apiv1.PaymentCustomer) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *customer) ToUpdate(r *apiv1.PaymentCustomer) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *customer) Update(rq any) (*apiv1.PaymentCustomer, error) {
	panic("unimplemented")
}
