package v1

import (
	"github.com/bufbuild/connect-go"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ip struct {
	c *config.Config
}

func NewIPCmd(c *config.Config) *cobra.Command {
	w := &ip{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*connect.Request[apiv1.IPServiceAllocateRequest], *connect.Request[apiv1.IPServiceStaticRequest], *apiv1.IP]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*connect.Request[apiv1.IPServiceAllocateRequest], *connect.Request[apiv1.IPServiceStaticRequest], *apiv1.IP](w).WithFS(c.Fs),
		Singular:        "ip",
		Plural:          "ips",
		Description:     "a ip address of metal-stack cloud",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project from where ips should be listed")
			genericcli.Must(cmd.MarkFlagRequired("project"))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project where the ip should be created")
			cmd.Flags().StringP("name", "", "", "name of the ip")
			cmd.Flags().StringP("description", "", "", "description of the ip")
			cmd.Flags().StringSliceP("tags", "", nil, "tags to add to the ip")
			cmd.Flags().BoolP("static", "", false, "make this ip static")
			cmd.MarkFlagsMutuallyExclusive("project", "file")
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "uuid of the ip")
			cmd.Flags().StringP("project", "", "", "project from where the ip should be made static")
			genericcli.Must(cmd.MarkFlagRequired("uuid"))
			genericcli.Must(cmd.MarkFlagRequired("project"))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project of the ip")
			genericcli.Must(cmd.MarkFlagRequired("project"))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project of the ip")
			genericcli.Must(cmd.MarkFlagRequired("project"))
		},
		CreateRequestFromCLI: func() (*connect.Request[apiv1.IPServiceAllocateRequest], error) {
			ipar := &apiv1.IPServiceAllocateRequest{
				Project:     viper.GetString("project"),
				Name:        viper.GetString("name"),
				Description: viper.GetString("description"),
				Tags:        viper.GetStringSlice("tags"),
				Static:      viper.GetBool("static"),
			}
			return connect.NewRequest(ipar), nil
		},
		UpdateRequestFromCLI: func(args []string) (*connect.Request[apiv1.IPServiceStaticRequest], error) {
			ipsr := &apiv1.IPServiceStaticRequest{
				Uuid:    viper.GetString("uuid"),
				Project: viper.GetString("project"),
			}
			return connect.NewRequest(ipsr), nil
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

// Create implements genericcli.CRUD
func (c *ip) Create(rq *connect.Request[apiv1.IPServiceAllocateRequest]) (*apiv1.IP, error) {
	resp, err := c.c.Apiv1Client.IP().Allocate(c.c.Ctx, rq)
	if err != nil {
		return nil, err
	}
	return resp.Msg.Ip, nil
}

// Delete implements genericcli.CRUD
func (c *ip) Delete(id string) (*apiv1.IP, error) {
	resp, err := c.c.Apiv1Client.IP().Delete(c.c.Ctx, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
		Project: viper.GetString("project"),
		Uuid:    id,
	}))
	if err != nil {
		return nil, err
	}
	return resp.Msg.Ip, nil
}

// Get implements genericcli.CRUD
func (c *ip) Get(id string) (*apiv1.IP, error) {
	resp, err := c.c.Apiv1Client.IP().Get(c.c.Ctx, connect.NewRequest(&apiv1.IPServiceGetRequest{
		Project: viper.GetString("project"),
		Uuid:    id,
	}))
	if err != nil {
		return nil, err
	}
	return resp.Msg.Ip, nil
}

// List implements genericcli.CRUD
func (c *ip) List() ([]*apiv1.IP, error) {
	// FIXME implement filters and paging
	resp, err := c.c.Apiv1Client.IP().List(c.c.Ctx, connect.NewRequest(&apiv1.IPServiceListRequest{
		Project: viper.GetString("project"),
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ips, nil
}

// ToCreate implements genericcli.CRUD
func (c *ip) ToCreate(r *apiv1.IP) (*connect.Request[apiv1.IPServiceAllocateRequest], error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *ip) ToUpdate(r *apiv1.IP) (*connect.Request[apiv1.IPServiceStaticRequest], error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *ip) Update(rq *connect.Request[apiv1.IPServiceStaticRequest]) (*apiv1.IP, error) {
	resp, err := c.c.Apiv1Client.IP().Static(c.c.Ctx, rq)
	if err != nil {
		return nil, err
	}
	return resp.Msg.Ip, nil
}
