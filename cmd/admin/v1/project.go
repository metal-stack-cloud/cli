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
		Description:     "a project in metalstack.cloud",
		Sorter:          sorters.ProjectSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().Uint64("limit", 100, "limit results returned")
			cmd.Flags().StringP("owner", "w", "", "filter by project owner")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *project) Create(rq any) (*apiv1.Project, error) {
	panic("unimplemented")
}

func (c *project) Delete(id string) (*apiv1.Project, error) {
	panic("unimplemented")
}

func (c *project) Get(id string) (*apiv1.Project, error) {
	panic("unimplemented")
}

var nextPage *uint64 

func (p *project) List() ([]*apiv1.Project, error) {
	ctx, cancel := p.c.NewRequestContext()
	defer cancel()

	req := &adminv1.ProjectServiceListRequest{}

	if viper.IsSet("limit") {
		req.Paging = &apiv1.Paging{
			Count: pointer.Pointer(viper.GetUint64("limit")),
			Page:  nextPage,
		}
	}
	

	resp, err := p.c.Client.Adminv1().Project().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	nextPage = resp.Msg.NextPage
	if nextPage != nil {
		err = p.c.ListPrinter.Print(resp.Msg.Projects)
		if err != nil {
			return nil, err
		}
		return p.List()
	}

	return resp.Msg.Projects, nil
}


func (c *project) Convert(r *apiv1.Project) (string, any, any, error) {
	panic("unimplemented")
}

func (c *project) Update(rq any) (*apiv1.Project, error) {
	panic("unimplemented")
}
