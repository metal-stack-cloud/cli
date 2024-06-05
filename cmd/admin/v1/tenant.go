package v1

import (
	"fmt"
	"strings"

	"connectrpc.com/connect"
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
		Description:     "a tenant of metalstack.cloud",
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
			ctx, cancel := c.NewRequestContext()
			defer cancel()

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
			resp, err := c.Client.Adminv1().Tenant().Admit(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to admit tenant: %w", err)
			}

			return c.DescribePrinter.Print(resp.Msg.Tenant)
		},
		ValidArgsFunction: c.Completion.AdminTenantListCompletion,
	}

	admitCmd.Flags().StringP("coupon-id", "", "", "optional add a coupon with given id, see coupon list for available coupons")

	genericcli.Must(admitCmd.RegisterFlagCompletionFunc("coupon-id", c.Completion.AdminPaymentCouponListCompletion))

	revokeCmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke a tenant",
		Long:  "revoke a tenant to be able to consume resources, can be enabled again with admit",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			req := &adminv1.TenantServiceRevokeRequest{
				TenantId: id,
			}
			resp, err := c.Client.Adminv1().Tenant().Revoke(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to revoke tenant: %w", err)
			}

			return c.DescribePrinter.Print(resp.Msg.Tenant)
		},
		ValidArgsFunction: c.Completion.AdminTenantListCompletion,
	}

	addMemberCmd := &cobra.Command{
		Use:   "add-member",
		Short: "Add a new member to a tenant",
		Long:  `Add a new member to an existing tenant by specifying the tenant ID, member's ID, and role.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("tenant-id", cmd.Flags().Lookup("tenant-id"))
			viper.BindPFlag("member-id", cmd.Flags().Lookup("member-id"))
			viper.BindPFlag("role", cmd.Flags().Lookup("role"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			tenantId := viper.GetString("tenant-id")
			memberId := viper.GetString("member-id")
			memberRole := viper.GetString("role")

			if tenantId == "" || memberId == "" || memberRole == "" {
				return fmt.Errorf("tenant ID, member ID, and role must all be specified")
			}

			roleValue, ok := apiv1.TenantRole_value[memberRole]
			if !ok {
				validRoles := getValidRoles()
				return fmt.Errorf("invalid role specified: %s. Valid roles are: %s", memberRole, strings.Join(validRoles, ", "))
			}

			req := &adminv1.TenantServiceAddMemberRequest{
				TenantId: tenantId,
				MemberId: memberId,
				Role:     apiv1.TenantRole(roleValue),
			}

			resp, err := c.Client.Adminv1().Tenant().AddMember(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to add member to tenant: %w", err)
			}

			return c.DescribePrinter.Print(resp.Msg.Member)
		},
	}

	addMemberCmd.Flags().String("tenant-id", "", "ID of the tenant where the member is added")
	addMemberCmd.Flags().String("member-id", "", "ID of the member to be added")
	addMemberCmd.Flags().String("role", "", "Role of the member within the tenant")

	genericcli.Must(addMemberCmd.MarkFlagRequired("tenant-id"))
	genericcli.Must(addMemberCmd.MarkFlagRequired("member-id"))
	genericcli.Must(addMemberCmd.MarkFlagRequired("role"))

	removeMemberCmd := &cobra.Command{
		Use:   "remove-member",
		Short: "Remove a member from a tenant",
		Long:  `Remove an existing member from a tenant by specifying the tenant ID and the member's ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			tenantId := viper.GetString("tenant-id")
			memberId := viper.GetString("member-id")

			if tenantId == "" || memberId == "" {
				return fmt.Errorf("both tenant ID and member ID must be specified")
			}

			req := &adminv1.TenantServiceRemoveMemberRequest{
				TenantId: tenantId,
				MemberId: memberId,
			}

			_, err := c.Client.Adminv1().Tenant().RemoveMember(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to remove member from tenant: %w", err)
			}

			return nil
		},
	}

	removeMemberCmd.Flags().String("tenant-id", "", "ID of the tenant from which the member is to be removed")
	removeMemberCmd.Flags().String("member-id", "", "ID of the member to remove")
	genericcli.Must(removeMemberCmd.MarkFlagRequired("tenant-id"))
	genericcli.Must(removeMemberCmd.MarkFlagRequired("member-id"))

	return genericcli.NewCmds(cmdsConfig, admitCmd, revokeCmd, addMemberCmd, removeMemberCmd)

}

func getValidRoles() []string {
	validRoles := make([]string, 0, len(apiv1.TenantRole_value))
	for r := range apiv1.TenantRole_value {
		validRoles = append(validRoles, r)
	}
	return validRoles
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
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	// TODO: implement filters and paging

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

	resp, err := c.c.Client.Adminv1().Tenant().List(ctx, connect.NewRequest(req))
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

func (c *tenant) Convert(r *apiv1.Tenant) (string, any, any, error) {
	panic("unimplemented")
}

func (c *tenant) Update(rq any) (*apiv1.Tenant, error) {
	panic("unimplemented")
}
