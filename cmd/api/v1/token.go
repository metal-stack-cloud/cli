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
			var permissions []*apiv1.MethodPermission
			for _, r := range viper.GetStringSlice("permissions") {
				project, semicolonSeparatedMethods, ok := strings.Cut(r, "=")
				if !ok {
					return nil, fmt.Errorf("permissions must be provided in the form <project>=<methods-colon-separated>")
				}

				permissions = append(permissions, &apiv1.MethodPermission{
					Subject: project,
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

			genericcli.Must(cmd.RegisterFlagCompletionFunc("permissions", c.Completion.TokenPermissionsCompletionfunc))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("roles", c.Completion.TokenRolesCompletion))

		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Aliases = append(cmd.Aliases, "revoke")
		},
		OnlyCmds:    genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DeleteCmd, genericcli.CreateCmd),
		ValidArgsFn: w.c.Completion.TokenListCompletion,
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (c *token) Get(id string) (*apiv1.Token, error) {
	panic("unimplemented")
}

func (c *token) List() ([]*apiv1.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.TokenServiceListRequest{}

	resp, err := c.c.Client.Apiv1().Token().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	return resp.Msg.GetTokens(), nil
}

func (c *token) Create(rq *apiv1.TokenServiceCreateRequest) (*apiv1.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Token().Create(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(c.c.Out, "Make sure to copy your personal access token now as you will not be able to see this again.\n")
	fmt.Fprintln(c.c.Out)
	fmt.Fprintln(c.c.Out, resp.Msg.GetSecret())
	fmt.Fprintln(c.c.Out)

	// TODO: allow printer in metal-lib to be silenced

	return resp.Msg.GetToken(), nil
}

func (c *token) Delete(id string) (*apiv1.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.TokenServiceRevokeRequest{
		Uuid: id,
	}

	_, err := c.c.Client.Apiv1().Token().Revoke(ctx, connect.NewRequest(req))
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
