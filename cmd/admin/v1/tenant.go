package v1

import (
	"fmt"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type tenant struct {
	c *config.Config
}

func newTenantCmd(c *config.Config) *cobra.Command {
	w := &tenant{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Tenant]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv1.Tenant](w).WithFS(c.Fs),
		Singular:        "tenant",
		Plural:          "tenants",
		Description:     "a tenant of metal-stack cloud",
		Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.Pf.NewPrinterDefaultYAML(c.Out) },
		ListPrinter:     func() printers.Printer { return c.Pf.NewPrinter(c.Out) },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
	}

	admitCmd := &cobra.Command{
		Use:   "admit",
		Short: "admit a tenant",
		Long:  "only admitted tenants are allowed to consume resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			resp, err := c.Adminv1Client.Tenant().Admit(c.Ctx, &adminv1.TenantServiceAdmitRequest{
				TenantId: id,
			})
			if err != nil {
				return fmt.Errorf("failed to admit tenant: %w", err)
			}

			return c.Pf.NewPrinter(c.Out).Print(resp.Tenant)
		},
	}

	return genericcli.NewCmds(cmdsConfig, admitCmd)
}

// Create implements genericcli.CRUD
func (c *tenant) Create(rq any) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (c *tenant) Delete(id string) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD
func (c *tenant) Get(id string) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

// List implements genericcli.CRUD
func (c *tenant) List() ([]*apiv1.Tenant, error) {
	// FIXME implement filters and paging
	resp, err := c.c.Adminv1Client.Tenant().List(c.c.Ctx, &adminv1.TenantServiceListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get tenants: %w", err)
	}

	return resp.Tenants, nil
}

// ToCreate implements genericcli.CRUD
func (c *tenant) ToCreate(r *apiv1.Tenant) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *tenant) ToUpdate(r *apiv1.Tenant) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *tenant) Update(rq any) (*apiv1.Tenant, error) {
	panic("unimplemented")
}
