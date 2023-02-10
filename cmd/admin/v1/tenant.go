package v1

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().BoolP("admitted", "a", false, "filter by admitted tenant")
			cmd.Flags().Uint64("limit", 100, "limit results returned")
			cmd.Flags().StringP("provider", "", "", "filter by provider")
			cmd.Flags().StringP("email", "", "", "filter by email")
		},
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
			req := &adminv1.TenantServiceAdmitRequest{
				TenantId: id,
			}
			if viper.IsSet("coupon-id") {
				req.CouponId = pointer.Pointer(viper.GetString("coupon-id"))
			}
			resp, err := c.Adminv1Client.Tenant().Admit(c.Ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to admit tenant: %w", err)
			}

			return c.DescribePrinter.Print(resp.Msg.Tenant)
		},
	}
	admitCmd.Flags().StringP("coupon-id", "", "", "optional add a coupon with given id, see coupon list for available coupons")

	revokeCmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke a tenant",
		Long:  "revoke a tenant to be able to consume resources, can be enabled again with admit",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			req := &adminv1.TenantServiceRevokeRequest{
				TenantId: id,
			}
			resp, err := c.Adminv1Client.Tenant().Revoke(c.Ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to revoke tenant: %w", err)
			}

			return c.DescribePrinter.Print(resp.Msg.Tenant)
		},
	}

	return genericcli.NewCmds(cmdsConfig, admitCmd, revokeCmd)
}

func (c *tenant) Create(rq any) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

func (c *tenant) Delete(id string) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

func (c *tenant) Get(id string) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

var nextpage *uint64

func (c *tenant) List() ([]*apiv1.Tenant, error) {
	// FIXME implement filters and paging

	req := &adminv1.TenantServiceListRequest{}
	if viper.IsSet("admitted") {
		req.Admitted = pointer.Pointer(viper.GetBool("admitted"))
	}
	if viper.IsSet("limit") {
		req.Paging = &apiv1.Paging{
			Count: pointer.Pointer(viper.GetUint64("limit")),
			Page:  nextpage,
		}
	}
	if viper.IsSet("provider") {
		return nil, fmt.Errorf("unimplemented filter by provider")
	}
	if viper.IsSet("email") {
		return nil, fmt.Errorf("unimplemented filter by provider")
	}
	resp, err := c.c.Adminv1Client.Tenant().List(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get tenants: %w", err)
	}

	nextpage = resp.Msg.NextPage
	if nextpage != nil {
		err = c.c.ListPrinter.Print(resp.Msg.Tenants)
		if err != nil {
			return nil, err
		}
		err := genericcli.PromptCustom(&genericcli.PromptConfig{
			Message:         "show more",
			No:              "n",
			AcceptedAnswers: genericcli.PromptDefaultAnswers(),
			ShowAnswers:     true,
		})
		if err != nil {
			return resp.Msg.Tenants, err
		}
		return c.List()
	}
	return resp.Msg.Tenants, nil
}

func (c *tenant) Update(rq any) (*apiv1.Tenant, error) {
	panic("unimplemented")
}

func (c *tenant) Convert(r *apiv1.Tenant) (string, any, any, error) {
	panic("unimplemented")
}
