package v1

import (
	"fmt"

	"connectrpc.com/connect"
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

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Project]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv1.Project](w).WithFS(c.Fs),
		Singular:        "project",
		Plural:          "projects",
		Description:     "manage api projects",
		Sorter:          sorters.ProjectSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "lists only projects with the given name")
			cmd.Flags().String("tenant", "", "lists only project with the given tenant")
		},
		ValidArgsFn: w.c.Completion.ProjectListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
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

func (t *project) Create(rq any) (*apiv1.Project, error) {
	panic("unimplemented")
}

func (t *project) Delete(id string) (*apiv1.Project, error) {
	panic("unimplemented")
}

func (t *project) Convert(r *apiv1.Project) (string, any, any, error) {
	panic("unimplemented")
}

func (t *project) Update(rq any) (*apiv1.Project, error) {
	panic("unimplemented")
}
