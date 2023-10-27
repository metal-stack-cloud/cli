package v1

import (
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type token struct {
	c *config.Config
}

func newTokenCmd(c *config.Config) *cobra.Command {
	w := &token{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.TokenServiceCreateRequest, any, *apiv1.Token]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*apiv1.TokenServiceCreateRequest, any, *apiv1.Token](w).WithFS(c.Fs),
		Singular:        "token",
		Plural:          "tokens",
		Description:     "manage api tokens for accessing the metalstack.cloud api",
		Sorter:          sorters.TokenSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		CreateRequestFromCLI: func() (*apiv1.TokenServiceCreateRequest, error) {
			var permissions []*apiv1.ProjectPermission
			for _, r := range viper.GetStringSlice("permissions") {
				project, semicolonSeparatedMethods, ok := strings.Cut(r, "=")
				if !ok {
					return nil, fmt.Errorf("permissions must be provided in the form <project>=<methods-colon-separated>")
				}

				permissions = append(permissions, &apiv1.ProjectPermission{
					Project: project,
					Methods: strings.Split(semicolonSeparatedMethods, ":"),
				})
			}

			var roles []*apiv1.TokenRole
			for _, r := range viper.GetStringSlice("roles") {
				subject, role, ok := strings.Cut(r, "=")
				if !ok {
					return nil, fmt.Errorf("roles must be provided in the form <subject>=<role>")
				}

				roles = append(roles, &apiv1.TokenRole{
					Subject: subject,
					Role:    role,
				})
			}

			return &apiv1.TokenServiceCreateRequest{
				// TODO: api should have an endpoint to list possible permissions and roles
				Description: viper.GetString("description"),
				Permissions: permissions,
				Roles:       roles,
				Expires:     durationpb.New(viper.GetDuration("expires")),
			}, nil
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("description", "", "a short description for the intention to use this token for")
			cmd.Flags().StringSlice("permissions", nil, "the permissions to associate with the api token in the form <project>=<methods-colon-separated>")
			cmd.Flags().StringSlice("roles", nil, "the roles to associate with the api token in the form <subject>=<role>")
			cmd.Flags().Duration("expires", 8*time.Hour, "the duration how long the api token is valid")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("permissions", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				methods, err := c.Client.Apiv1().Method().TokenScopedList(c.Ctx, connect.NewRequest(&apiv1.MethodServiceTokenScopedListRequest{}))
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				subject := ""
				if s, _, ok := strings.Cut(toComplete, "="); ok {
					subject = s
				}

				if subject == "" {
					var perms []string

					for _, p := range methods.Msg.Permissions {
						perms = append(perms, p.Project)
					}

					return perms, cobra.ShellCompDirectiveNoFileComp
				}

				// FIXME: completion does not work at this point, investigate why

				var perms []string

				for _, p := range methods.Msg.Permissions {
					for _, m := range p.Methods {
						perms = append(perms, m)
					}
				}

				return perms, cobra.ShellCompDirectiveDefault
			}))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("roles", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				methods, err := c.Client.Apiv1().Method().TokenScopedList(c.Ctx, connect.NewRequest(&apiv1.MethodServiceTokenScopedListRequest{}))
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				var roles []string

				for _, r := range methods.Msg.Roles {
					roles = append(roles, r.Subject+":"+r.Role)
				}

				return roles, cobra.ShellCompDirectiveNoFileComp
			}))

		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Aliases = append(cmd.Aliases, "revoke")
		},
		OnlyCmds:    genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DeleteCmd, genericcli.CreateCmd),
		ValidArgsFn: w.c.Completion.TokenListCompletion,
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (t *token) Get(id string) (*apiv1.Token, error) {
	panic("unimplemented")
}

func (t *token) List() ([]*apiv1.Token, error) {
	req := &apiv1.TokenServiceListRequest{}

	resp, err := t.c.Client.Apiv1().Token().List(t.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	return resp.Msg.GetTokens(), nil
}

func (t *token) Create(rq *apiv1.TokenServiceCreateRequest) (*apiv1.Token, error) {
	resp, err := t.c.Client.Apiv1().Token().Create(t.c.Ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(t.c.Out, "Make sure to copy your personal access token now as you will not be able to see this again.\n")
	fmt.Fprintln(t.c.Out)
	fmt.Fprintln(t.c.Out, resp.Msg.GetSecret())
	fmt.Fprintln(t.c.Out)

	// TODO: allow printer in metal-lib to be silenced

	return resp.Msg.GetToken(), nil
}

func (t *token) Delete(id string) (*apiv1.Token, error) {
	req := &apiv1.TokenServiceRevokeRequest{
		Uuid: id,
	}

	_, err := t.c.Client.Apiv1().Token().Revoke(t.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	return &apiv1.Token{
		Uuid: id,
	}, nil
}

func (t *token) Convert(r *apiv1.Token) (string, *apiv1.TokenServiceCreateRequest, any, error) {
	return r.Uuid, &apiv1.TokenServiceCreateRequest{
		Description: r.GetDescription(),
		Permissions: r.GetPermissions(),
		Roles:       r.GetRoles(),
		Expires:     durationpb.New(time.Until(r.GetExpires().AsTime())),
	}, nil, nil
}

func (t *token) Update(rq any) (*apiv1.Token, error) {
	panic("unimplemented")
}
