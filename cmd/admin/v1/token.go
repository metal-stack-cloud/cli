package v1

import (
	"fmt"

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

type token struct {
	c *config.Config
}

func newTokenCmd(c *config.Config) *cobra.Command {
	w := &token{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Token]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv1.Token](w).WithFS(c.Fs),
		Singular:        "token",
		Plural:          "tokens",
		Description:     "manage api tokens for accessing the metalstack.cloud api",
		Sorter:          sorters.TokenSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("user", "", "the uuid of the user to list the tokens for")
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("user", "", "the uuid of the user who owns the token")

			cmd.Aliases = append(cmd.Aliases, "revoke")
		},
		OnlyCmds:    genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DeleteCmd),
		ValidArgsFn: w.c.Completion.TokenListCompletion,
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (t *token) Get(id string) (*apiv1.Token, error) {
	panic("unimplemented")
}

func (c *token) List() ([]*apiv1.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv1.TokenServiceListRequest{}

	if viper.IsSet("user") {
		req.UserId = pointer.Pointer(viper.GetString("user"))
	}

	resp, err := c.c.Client.Adminv1().Token().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	return resp.Msg.GetTokens(), nil
}

func (t *token) Create(rq any) (*apiv1.Token, error) {
	panic("unimplemented")
}

func (c *token) Delete(id string) (*apiv1.Token, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	if !viper.IsSet("user") {
		return nil, fmt.Errorf("user is required to be set")
	}

	req := &adminv1.TokenServiceRevokeRequest{
		Uuid:   id,
		UserId: viper.GetString("user"),
	}

	_, err := c.c.Client.Adminv1().Token().Revoke(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to revoke token: %w", err)
	}

	return &apiv1.Token{
		Uuid: id,
	}, nil
}

func (t *token) Convert(r *apiv1.Token) (string, any, any, error) {
	panic("unimplemented")
}

func (t *token) Update(rq any) (*apiv1.Token, error) {
	panic("unimplemented")
}
