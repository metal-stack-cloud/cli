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
	c *config.Config
}

func newCouponCmd(c *config.Config) *cobra.Command {
	w := &coupon{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Coupon]{
		BinaryName:  config.BinaryName,
		GenericCLI:  genericcli.NewGenericCLI[any, any, *apiv1.Coupon](w).WithFS(c.Fs),
		Singular:    "coupon",
		Plural:      "coupons",
		Description: "coupon related actions of metalstack.cloud",
		// Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
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
	resp, err := c.c.Client.Adminv1().Payment().ListCoupons(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get coupons: %w", err)
	}
	return resp.Msg.Coupons, nil
}

// Convert implements genericcli.CRUD
func (c *coupon) Convert(r *apiv1.Coupon) (string, any, any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *coupon) Update(rq any) (*apiv1.Coupon, error) {
	panic("unimplemented")
}
