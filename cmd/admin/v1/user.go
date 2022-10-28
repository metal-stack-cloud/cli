package v1

import (
	"fmt"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type user struct {
	c *config.Config
}

func NewUserCmd(c *config.Config) *cobra.Command {
	w := &user{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *adminv1.User]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *adminv1.User](w).WithFS(c.Fs),
		Singular:        "user",
		Plural:          "users",
		Description:     "a user of metal-stack cloud",
		Sorter:          sorters.UserSorter(),
		DescribePrinter: func() printers.Printer { return c.Pf.NewPrinterDefaultYAML(c.Out) },
		ListPrinter:     func() printers.Printer { return c.Pf.NewPrinter(c.Out) },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
	}

	admitCmd := &cobra.Command{
		Use:   "admit",
		Short: "admit a user",
		Long:  "only admitted users are allowed to consume resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			resp, err := c.Adminv1Client.User().Admit(c.Ctx, &adminv1.UserServiceAdmitRequest{
				UserId: id,
			})
			if err != nil {
				return fmt.Errorf("failed to admit user: %w", err)
			}

			return c.Pf.NewPrinter(c.Out).Print(resp.User)
		},
	}

	return genericcli.NewCmds(cmdsConfig, admitCmd)
}

// Create implements genericcli.CRUD
func (c *user) Create(rq any) (*adminv1.User, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (c *user) Delete(id string) (*adminv1.User, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD
func (c *user) Get(id string) (*adminv1.User, error) {
	panic("unimplemented")
}

// List implements genericcli.CRUD
func (c *user) List() ([]*adminv1.User, error) {
	resp, err := c.c.Adminv1Client.User().List(c.c.Ctx, &adminv1.UserServiceListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return resp.Users, nil
}

// ToCreate implements genericcli.CRUD
func (c *user) ToCreate(r *adminv1.User) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *user) ToUpdate(r *adminv1.User) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *user) Update(rq any) (*adminv1.User, error) {
	panic("unimplemented")
}
