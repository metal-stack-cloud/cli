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

type project struct {
	c *config.Config
}

func newProjectCmd(c *config.Config) *cobra.Command {
	w := &project{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.ProjectServiceCreateRequest, *apiv1.ProjectServiceUpdateRequest, *apiv1.Project]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*apiv1.ProjectServiceCreateRequest, *apiv1.ProjectServiceUpdateRequest, *apiv1.Project](w).WithFS(c.Fs),
		Singular:        "project",
		Plural:          "projects",
		Description:     "manage api projects",
		Sorter:          sorters.ProjectSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "lists only projects with the given name")
			cmd.Flags().String("tenant", "", "lists only project with the given tenant")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the project to create")
			cmd.Flags().String("description", "", "the description of the project to create")
			cmd.Flags().String("tenant", "", "the tenant of this project, defaults to tenant of the default project")
		},
		CreateRequestFromCLI: w.createRequestFromCLI,
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "the name of the project to update")
			cmd.Flags().String("description", "", "the description of the project to update")
		},
		UpdateRequestFromCLI: w.updateRequestFromCLI,
		ValidArgsFn:          w.c.Completion.ProjectListCompletion,
	}

	inviteCmd := &cobra.Command{
		Use:   "invite",
		Short: "manage project invites",
	}

	generateInviteCmd := &cobra.Command{
		Use:   "generate-join-secret",
		Short: "generate an invite secret to share with the new member",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.generateInvite()
		},
	}

	generateInviteCmd.Flags().StringP("project", "p", "", "the project for which to generate the invite")
	generateInviteCmd.Flags().String("role", apiv1.ProjectRole_PROJECT_ROLE_VIEWER.String(), "the role that the new member will assume when joining through the invite secret")

	genericcli.Must(generateInviteCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
	genericcli.Must(generateInviteCmd.RegisterFlagCompletionFunc("role", c.Completion.ProjectRoleCompletion))

	deleteInviteCmd := &cobra.Command{
		Use:     "delete <secret>",
		Aliases: []string{"destroy", "rm", "remove"},
		Short:   "deletes a pending invite",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.deleteInvite(args)
		},
		ValidArgsFunction: c.Completion.ProjectInviteListCompletion,
	}

	deleteInviteCmd.Flags().StringP("project", "p", "", "the project in which to delete the invite")

	genericcli.Must(deleteInviteCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

	listInvitesCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "lists the currently pending invites",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.listInvites()
		},
	}

	listInvitesCmd.Flags().StringP("project", "p", "", "the project for which to list the invites")

	genericcli.AddSortFlag(listInvitesCmd, sorters.ProjectInviteSorter())

	genericcli.Must(listInvitesCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

	joinProjectCmd := &cobra.Command{
		Use:   "join <secret>",
		Short: "join a project of someone who shared an invite secret with you",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.join(args)
		},
	}

	inviteCmd.AddCommand(generateInviteCmd, deleteInviteCmd, listInvitesCmd, joinProjectCmd)

	return genericcli.NewCmds(cmdsConfig, joinProjectCmd, inviteCmd)
}

func (c *project) Get(id string) (*apiv1.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.ProjectServiceGetRequest{
		Project: id,
	}

	resp, err := c.c.Client.Apiv1().Project().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return resp.Msg.GetProject(), nil
}

func (c *project) List() ([]*apiv1.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.ProjectServiceListRequest{
		Name:   pointer.PointerOrNil(viper.GetString("name")),
		Tenant: pointer.PointerOrNil(viper.GetString("tenant")),
	}

	resp, err := c.c.Client.Apiv1().Project().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return resp.Msg.GetProjects(), nil
}

func (c *project) Create(rq *apiv1.ProjectServiceCreateRequest) (*apiv1.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Project().Create(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return resp.Msg.Project, nil
}

func (c *project) Delete(id string) (*apiv1.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Project().Delete(ctx, connect.NewRequest(&apiv1.ProjectServiceDeleteRequest{
		Project: id,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to delete project: %w", err)
	}

	return resp.Msg.Project, nil
}

func (c *project) Convert(r *apiv1.Project) (string, *apiv1.ProjectServiceCreateRequest, *apiv1.ProjectServiceUpdateRequest, error) {
	return r.Uuid, &apiv1.ProjectServiceCreateRequest{
			Login:       r.Tenant,
			Name:        r.Name,
			Description: r.Description,
		}, &apiv1.ProjectServiceUpdateRequest{
			Project:     r.Uuid,
			Name:        pointer.Pointer(r.Name),
			Description: pointer.Pointer(r.Description),
		}, nil
}

func (c *project) Update(rq *apiv1.ProjectServiceUpdateRequest) (*apiv1.Project, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Project().Update(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return resp.Msg.Project, nil
}

func (c *project) createRequestFromCLI() (*apiv1.ProjectServiceCreateRequest, error) {
	tenant := viper.GetString("tenant")

	if tenant == "" && c.c.GetProject() != "" {
		project, err := c.Get(c.c.GetProject())
		if err != nil {
			return nil, fmt.Errorf("unable to derive tenant from project: %w", err)
		}

		tenant = project.Tenant
	}

	if viper.GetString("name") == "" {
		return nil, fmt.Errorf("name must be given")
	}
	if viper.GetString("description") == "" {
		return nil, fmt.Errorf("description must be given")
	}

	return &apiv1.ProjectServiceCreateRequest{
		Login:       tenant,
		Name:        viper.GetString("name"),
		Description: viper.GetString("description"),
	}, nil
}

func (c *project) updateRequestFromCLI(args []string) (*apiv1.ProjectServiceUpdateRequest, error) {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	return &apiv1.ProjectServiceUpdateRequest{
		Project:     id,
		Name:        pointer.PointerOrNil(viper.GetString("name")),
		Description: pointer.PointerOrNil(viper.GetString("description")),
	}, nil
}

func (c *project) join(args []string) error {
	secret, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Project().InviteAccept(ctx, connect.NewRequest(&apiv1.ProjectServiceInviteAcceptRequest{
		Secret: secret,
	}))
	if err != nil {
		return fmt.Errorf("failed to join project: %w", err)
	}

	fmt.Fprintf(c.c.Out, "%s successfully joined project \"%s\"\n", color.GreenString("✔"), color.GreenString(resp.Msg.ProjectName))

	return nil
}

func (c *project) generateInvite() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Project().Invite(ctx, connect.NewRequest(&apiv1.ProjectServiceInviteRequest{
		Project: c.c.GetProject(),
		Role:    apiv1.ProjectRole(apiv1.ProjectRole_value[viper.GetString("role")]),
	}))
	if err != nil {
		return fmt.Errorf("failed to generate an invite: %w", err)
	}

	fmt.Fprintf(c.c.Out, "You can share this secret with the member to join, it expires in %s:\n\n", humanize.Time(resp.Msg.Invite.ExpiresAt.AsTime()))
	fmt.Fprintf(c.c.Out, "%s (https://console.metalstack.cloud/invite/%s)\n", resp.Msg.Invite.Secret, resp.Msg.Invite.Secret)

	return nil
}

func (c *project) listInvites() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Project().InvitesList(ctx, connect.NewRequest(&apiv1.ProjectServiceInvitesListRequest{
		Project: c.c.GetProject(),
	}))
	if err != nil {
		return fmt.Errorf("failed to list invites: %w", err)
	}

	err = sorters.ProjectInviteSorter().SortBy(resp.Msg.Invites)
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(resp.Msg.Invites)
}

func (c *project) deleteInvite(args []string) error {
	secret, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	_, err = c.c.Client.Apiv1().Project().InviteDelete(ctx, connect.NewRequest(&apiv1.ProjectServiceInviteDeleteRequest{
		Project: c.c.GetProject(),
		Secret:  secret,
	}))
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}
