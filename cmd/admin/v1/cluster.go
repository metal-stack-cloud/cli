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
	"github.com/spf13/cobra"
)

type cluster struct {
	c *config.Config
}

func newClusterCmd(c *config.Config) *cobra.Command {
	w := &cluster{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Cluster]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv1.Cluster](w).WithFS(c.Fs),
		Singular:        "cluster",
		Plural:          "clusters",
		Description:     "cluster related actions of metalstack.cloud",
		Sorter:          sorters.ClusterSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
	}
	return genericcli.NewCmds(cmdsConfig)
}

// Create implements genericcli.CRUD
func (c *cluster) Create(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD
func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	req := &adminv1.ClusterServiceGetRequest{
		Uuid: id,
	}
	resp, err := c.c.Adminv1Client.Cluster().Get(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}
	return resp.Msg.Cluster, nil
}

// List implements genericcli.CRUD
func (c *cluster) List() ([]*apiv1.Cluster, error) {
	// FIXME implement filters and paging

	req := &adminv1.ClusterServiceListRequest{}
	resp, err := c.c.Adminv1Client.Cluster().List(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}
	return resp.Msg.Clusters, nil
}

// ToCreate implements genericcli.CRUD
func (c *cluster) ToCreate(r *apiv1.Cluster) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *cluster) ToUpdate(r *apiv1.Cluster) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *cluster) Update(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}
