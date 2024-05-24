package v1

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
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

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.TenantServiceCreateRequest, *apiv1.TenantServiceUpdateRequest, *apiv1.Tenant]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*apiv1.TenantServiceCreateRequest, *apiv1.TenantServiceUpdateRequest, *apiv1.Tenant](w).WithFS(c.Fs),
		Singular:        "tenant",
		Plural:          "tenants",
		Description:     "manage api tenants",
		Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "lists only tenants with the given name")
			cmd.Flags().String("id", "", "lists only tenant with the given tenant id")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant to create")
			cmd.Flags().String("description", "", "the description of the tenant to create")
			cmd.Flags().String("email", "", "the email of the tenant to create")
			cmd.Flags().String("phone", "", "the phone number of the tenant to create")
			cmd.Flags().String("avatar-url", "", "the avatar url of the tenant to create")
		},
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the tenant to update")
			cmd.Flags().String("description", "", "the description of the tenant to update")
		},
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		ValidArgsFn:          w.c.Completion.TenantListCompletion,
	}

	inviteCmd := &cobra.Command{
		Use:   "invite",
		Short: "manage tenant invites",
	}

	generateInviteCmd := &cobra.Command{
		Use:   "generate-join-secret",
		Short: "generate an invite secret to share with the new member",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.generateInvite()
		},
	}

	generateInviteCmd.Flags().String("tenant", "", "the tenant for which to generate the invite")
	generateInviteCmd.Flags().String("role", apiv1.TenantRole_TENANT_ROLE_VIEWER.String(), "the role that the new member will assume when joining through the invite secret")

	genericcli.Must(generateInviteCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
	genericcli.Must(generateInviteCmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRoleCompletion))

	deleteInviteCmd := &cobra.Command{
		Use:     "delete <secret>",
		Aliases: []string{"destroy", "rm", "remove"},
		Short:   "deletes a pending invite",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.deleteInvite(args)
		},
		ValidArgsFunction: c.Completion.TenantInviteListCompletion,
	}

	deleteInviteCmd.Flags().String("tenant", "", "the tenant in which to delete the invite")

	genericcli.Must(deleteInviteCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	listInvitesCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "lists the currently pending invites",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.listInvites()
		},
	}

	listInvitesCmd.Flags().String("tenant", "", "the tenant for which to list the invites")

	genericcli.AddSortFlag(listInvitesCmd, sorters.TenantInviteSorter())

	genericcli.Must(listInvitesCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	joinTenantCmd := &cobra.Command{
		Use:   "join <secret>",
		Short: "join a tenant of someone who shared an invite secret with you",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.join(args)
		},
	}

	memberCmd := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "manage tenant members",
	}

	listMembersCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "lists members of a tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.listMembers()
		},
	}

	listMembersCmd.Flags().String("tenant", "", "the tenant in which to remove the member")

	genericcli.AddSortFlag(listMembersCmd, sorters.TenantMemberSorter())

	removeMemberCmd := &cobra.Command{
		Use:     "remove <member>",
		Short:   "remove member from a tenant",
		Aliases: []string{"destroy", "rm", "remove"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.removeMember(args)
		},
		ValidArgsFunction: c.Completion.TenantMemberListCompletion,
	}

	removeMemberCmd.Flags().String("tenant", "", "the tenant in which to remove the member")

	genericcli.Must(removeMemberCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))

	updateMemberCmd := &cobra.Command{
		Use:   "update <member>",
		Short: "update member from a tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.updateMember(args)
		},
		ValidArgsFunction: c.Completion.TenantMemberListCompletion,
	}

	updateMemberCmd.Flags().String("tenant", "", "the tenant in which to remove the member")
	updateMemberCmd.Flags().String("role", "", "the role of the member")

	genericcli.Must(updateMemberCmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
	genericcli.Must(updateMemberCmd.RegisterFlagCompletionFunc("role", c.Completion.TenantRoleCompletion))

	memberCmd.AddCommand(removeMemberCmd, updateMemberCmd, listMembersCmd)

	inviteCmd.AddCommand(generateInviteCmd, deleteInviteCmd, listInvitesCmd, joinTenantCmd)

	return genericcli.NewCmds(cmdsConfig, joinTenantCmd, inviteCmd, memberCmd)
}

func (c *tenant) Get(id string) (*apiv1.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.TenantServiceGetRequest{
		Login: id,
	}

	resp, err := c.c.Client.Apiv1().Tenant().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return resp.Msg.GetTenant(), nil
}

func (c *tenant) List() ([]*apiv1.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.TenantServiceListRequest{
		Name: pointer.PointerOrNil(viper.GetString("name")),
		Id:   pointer.PointerOrNil(viper.GetString("tenant")),
	}
	resp, err := c.c.Client.Apiv1().Tenant().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return resp.Msg.GetTenants(), nil
}

func (c *tenant) Create(rq *apiv1.TenantServiceCreateRequest) (*apiv1.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Tenant().Create(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return resp.Msg.Tenant, nil
}

func (c *tenant) Delete(id string) (*apiv1.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Tenant().Delete(ctx, connect.NewRequest(&apiv1.TenantServiceDeleteRequest{
		Login: id,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to delete tenant: %w", err)
	}

	return resp.Msg.Tenant, nil
}

func (c *tenant) Convert(r *apiv1.Tenant) (string, *apiv1.TenantServiceCreateRequest, *apiv1.TenantServiceUpdateRequest, error) {
	var paymentDetails *apiv1.PaymentDetailsUpdate
	if r.PaymentDetails != nil {
		paymentDetails = &apiv1.PaymentDetailsUpdate{
			CustomerId:      pointer.PointerOrNil(r.PaymentDetails.CustomerId),
			PaymentMethodId: r.PaymentDetails.PaymentMethodId,
			SubscriptionId:  pointer.PointerOrNil(r.PaymentDetails.SubscriptionId),
			Vat:             pointer.PointerOrNil(r.PaymentDetails.Vat),
		}
	}

	return r.Login, &apiv1.TenantServiceCreateRequest{
			Name:        r.Name,
			Description: &r.Description,
			Email:       &r.Email,
			AvatarUrl:   &r.AvatarUrl,
			PhoneNumber: &r.PhoneNumber,
		},
		&apiv1.TenantServiceUpdateRequest{
			Login:          r.Login,
			Name:           pointer.PointerOrNil(r.Name),
			Email:          pointer.PointerOrNil(r.Email),
			AvatarUrl:      pointer.PointerOrNil(r.AvatarUrl),
			PaymentDetails: paymentDetails,
		},
		nil
}

func (c *tenant) Update(rq *apiv1.TenantServiceUpdateRequest) (*apiv1.Tenant, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Tenant().Update(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return resp.Msg.Tenant, nil
}

func (c *tenant) createRequestFromCLI() (*apiv1.TenantServiceCreateRequest, error) {
	return &apiv1.TenantServiceCreateRequest{
		Name:        viper.GetString("name"),
		Description: pointer.PointerOrNil(viper.GetString("description")),
		Email:       pointer.PointerOrNil(viper.GetString("email")),
		AvatarUrl:   pointer.PointerOrNil(viper.GetString("phone")),
		PhoneNumber: pointer.PointerOrNil(viper.GetString("avatar-url")),
	}, nil
}

func (c *tenant) updateRequestFromCLI(args []string) (*apiv1.TenantServiceUpdateRequest, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *tenant) join(args []string) error {
	secret, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Tenant().InviteGet(ctx, connect.NewRequest(&apiv1.TenantServiceInviteGetRequest{
		Secret: secret,
	}))
	if err != nil {
		return fmt.Errorf("failed to get tenant invite: %w", err)
	}

	err = genericcli.PromptCustom(&genericcli.PromptConfig{
		ShowAnswers: true,
		Message: fmt.Sprintf(
			"Do you want to join tenant \"%s\" as %s?",
			color.GreenString(resp.Msg.GetInvite().GetTargetTenantName()),
			resp.Msg.GetInvite().GetRole().String(),
		),
		In:  c.c.In,
		Out: c.c.Out,
	})
	if err != nil {
		return err
	}

	ctx2, cancel2 := c.c.NewRequestContext()
	defer cancel2()

	acceptResp, err := c.c.Client.Apiv1().Tenant().InviteAccept(ctx2, connect.NewRequest(&apiv1.TenantServiceInviteAcceptRequest{
		Secret: secret,
	}))
	if err != nil {
		return fmt.Errorf("failed to join tenant: %w", err)
	}

	fmt.Fprintf(c.c.Out, "%s successfully joined tenant \"%s\"\n", color.GreenString("✔"), color.GreenString(acceptResp.Msg.TenantName))

	return nil
}

func (c *tenant) generateInvite() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := c.c.Client.Apiv1().Tenant().Invite(ctx, connect.NewRequest(&apiv1.TenantServiceInviteRequest{
		Login: tenant,
		Role:  apiv1.TenantRole(apiv1.TenantRole_value[viper.GetString("role")]),
	}))
	if err != nil {
		return fmt.Errorf("failed to generate an invite: %w", err)
	}

	fmt.Fprintf(c.c.Out, "You can share this secret with the member to join, it expires in %s:\n\n", humanize.Time(resp.Msg.Invite.ExpiresAt.AsTime()))
	fmt.Fprintf(c.c.Out, "%s (https://console.metalstack.cloud/organization-invite/%s)\n", resp.Msg.Invite.Secret, resp.Msg.Invite.Secret)

	return nil
}

func (c *tenant) listInvites() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := c.c.Client.Apiv1().Tenant().InvitesList(ctx, connect.NewRequest(&apiv1.TenantServiceInvitesListRequest{
		Login: tenant,
	}))
	if err != nil {
		return fmt.Errorf("failed to list invites: %w", err)
	}

	err = sorters.TenantInviteSorter().SortBy(resp.Msg.Invites)
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.Msg.Invites)
}

func (c *tenant) deleteInvite(args []string) error {
	secret, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	_, err = c.c.Client.Apiv1().Tenant().InviteDelete(ctx, connect.NewRequest(&apiv1.TenantServiceInviteDeleteRequest{
		Login:  tenant,
		Secret: secret,
	}))
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

func (c *tenant) removeMember(args []string) error {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	_, err = c.c.Client.Apiv1().Tenant().RemoveMember(ctx, connect.NewRequest(&apiv1.TenantServiceRemoveMemberRequest{
		Login:    tenant,
		MemberId: member,
	}))
	if err != nil {
		return fmt.Errorf("failed to remove member from tenant: %w", err)
	}

	fmt.Fprintf(c.c.Out, "%s successfully removed member %q\n", color.GreenString("✔"), member)

	return nil
}

func (c *tenant) updateMember(args []string) error {
	member, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Tenant().UpdateMember(ctx, connect.NewRequest(&apiv1.TenantServiceUpdateMemberRequest{
		Login:    tenant,
		MemberId: member,
		Role:     apiv1.TenantRole(apiv1.TenantRole_value[viper.GetString("role")]),
	}))
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	return c.c.DescribePrinter.Print(resp.Msg.GetTenantMember())
}

func (c *tenant) listMembers() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	tenant, err := c.c.GetTenant()
	if err != nil {
		return err
	}

	resp, err := c.c.Client.Apiv1().Tenant().Get(ctx, connect.NewRequest(&apiv1.TenantServiceGetRequest{
		Login: tenant,
	}))
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	members := resp.Msg.GetTenantMembers()

	if err := sorters.TenantMemberSorter().SortBy(members); err != nil {
		return err
	}

	return c.c.ListPrinter.Print(members)
}
