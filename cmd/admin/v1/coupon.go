package v1

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type coupon struct {
	c               *config.Config
	listPrinter     func() printers.Printer
	describePrinter func() printers.Printer
}

func newCouponCmd(c *config.Config) *cobra.Command {
	w := &coupon{
		c: c,
	}
	w.listPrinter = func() printers.Printer { return c.Pf.NewPrinter(c.Out) }
	w.describePrinter = func() printers.Printer { return c.Pf.NewPrinterDefaultYAML(c.Out) }

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Coupon]{
		BinaryName:  config.BinaryName,
		GenericCLI:  genericcli.NewGenericCLI[any, any, *apiv1.Coupon](w).WithFS(c.Fs),
		Singular:    "coupon",
		Plural:      "coupons",
		Description: "coupon related actions of metal-stack cloud",
		// Sorter:          sorters.TenantSorter(),
		DescribePrinter: w.describePrinter,
		ListPrinter:     w.listPrinter,
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
	}
	return genericcli.NewCmds(cmdsConfig)
}

// Create implements genericcli.CRUD
func (c *coupon) Create(rq any) (*apiv1.Coupon, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (c *coupon) Delete(id string) (*apiv1.Coupon, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD
func (c *coupon) Get(id string) (*apiv1.Coupon, error) {
	panic("unimplemented")
}

// List implements genericcli.CRUD
func (c *coupon) List() ([]*apiv1.Coupon, error) {
	// FIXME implement filters and paging

	req := &adminv1.PaymentServiceListCouponsRequest{}
	resp, err := c.c.Adminv1Client.Payment().ListCoupons(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get coupons: %w", err)
	}
	return resp.Msg.Coupons, nil
}

// ToCreate implements genericcli.CRUD
func (c *coupon) ToCreate(r *apiv1.Coupon) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *coupon) ToUpdate(r *apiv1.Coupon) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *coupon) Update(rq any) (*apiv1.Coupon, error) {
	panic("unimplemented")
}
