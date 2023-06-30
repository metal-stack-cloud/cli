package v1

import (
	"fmt"

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

	cmdsConfig := &genericcli.CmdsConfig[*connect.Request[apiv1.IPServiceAllocateRequest], *apiv1.IP, *apiv1.IP]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*connect.Request[apiv1.IPServiceAllocateRequest], *apiv1.IP, *apiv1.IP](w).WithFS(c.Fs),
		Singular:        "ip",
		Plural:          "ips",
		Description:     "an ip address of metalstack.cloud",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project from where ips should be listed")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project where the ip should be created")
			cmd.Flags().StringP("name", "", "", "name of the ip")
			cmd.Flags().StringP("description", "", "", "description of the ip")
			cmd.Flags().StringSliceP("tags", "", nil, "tags to add to the ip")
			cmd.Flags().BoolP("static", "", false, "make this ip static")
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("uuid", "", "uuid of the ip")
			cmd.Flags().String("project", "", "project from where the ip should be made static")
			cmd.Flags().String("name", "", "name of the ip")
			cmd.Flags().String("description", "", "description of the ip")
			cmd.Flags().StringSlice("tags", nil, "tags of the ip")
			cmd.Flags().Bool("static", false, "make this ip static")
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project of the ip")
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project of the ip")
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
		UpdateRequestFromCLI: w.updateFromCLI,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) updateFromCLI(args []string) (*apiv1.IP, error) {
	ipToUpdate, err := c.Get(viper.GetString("uuid"))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve ip: %w", err)
	}

	if viper.IsSet("name") {
		ipToUpdate.Name = viper.GetString("name")
	}
	if viper.IsSet("description") {
		ipToUpdate.Description = viper.GetString("description")
	}
	if viper.IsSet("static") {
		ipToUpdate.Type = ipStaticToType(viper.GetBool("static"))
	}
	if viper.IsSet("tags") {
		ipToUpdate.Tags = viper.GetStringSlice("tags")
	}

	return ipToUpdate, nil
}



// Convert implements genericcli.CRUD.
func (*ip) Convert(r R) (string, C, U, error) {
	panic("unimplemented")
}

// Create implements genericcli.CRUD.
func (*ip) Create(rq C) (R, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD.
func (*ip) Delete(id string) (R, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD.
func (*ip) Get(id string) (R, error) {
	panic("unimplemented")
}

// List implements genericcli.CRUD.
func (*ip) List() ([]R, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD.
func (*ip) Update(rq U) (R, error) {
	panic("unimplemented")
}


// // Create implements genericcli.CRUD
// func (c *ip) Create(rq *connect.Request[apiv1.IPServiceAllocateRequest]) (*apiv1.IP, error) {
// 	resp, err := c.c.Apiv1Client.IP().Allocate(c.c.Ctx, rq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Msg.Ip, nil
// }

// // Delete implements genericcli.CRUD
// func (c *ip) Delete(id string) (*apiv1.IP, error) {
// 	project := viper.GetString("project")
// 	if project == "" {
// 		return nil, fmt.Errorf("project must be provided")
// 	}

// 	resp, err := c.c.Apiv1Client.IP().Delete(c.c.Ctx, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
// 		Project: project,
// 		Uuid:    id,
// 	}))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Msg.Ip, nil
// }

// // Get implements genericcli.CRUD
// func (c *ip) Get(id string) (*apiv1.IP, error) {
// 	project := viper.GetString("project")
// 	if project == "" {
// 		return nil, fmt.Errorf("project must be provided")
// 	}

// 	resp, err := c.c.Apiv1Client.IP().Get(c.c.Ctx, connect.NewRequest(&apiv1.IPServiceGetRequest{
// 		Project: project,
// 		Uuid:    id,
// 	}))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Msg.Ip, nil
// }

// // List implements genericcli.CRUD
// func (c *ip) List() ([]*apiv1.IP, error) {
// 	project := viper.GetString("project")
// 	if project == "" {
// 		return nil, fmt.Errorf("project must be provided")
// 	}

// 	// FIXME implement filters and paging
// 	resp, err := c.c.Apiv1Client.IP().List(c.c.Ctx, connect.NewRequest(&apiv1.IPServiceListRequest{
// 		Project: project,
// 	}))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp.Msg.Ips, nil
// }

// // ToCreate implements genericcli.CRUD
// func (c *ip) ToCreate(r *apiv1.IP) (*connect.Request[apiv1.IPServiceAllocateRequest], error) {
// 	return ipResponseToCreate(r), nil
// }

// // ToUpdate implements genericcli.CRUD
// func (c *ip) ToUpdate(r *apiv1.IP) (*apiv1.IP, error) {
// 	return ipResponseToUpdate(r), nil
// }

// // Update implements genericcli.CRUD
// func (c *ip) Update(rq *apiv1.IP) (*apiv1.IP, error) {
// 	resp, err := c.c.Apiv1Client.IP().Update(c.c.Ctx, &connect.Request[apiv1.IPServiceUpdateRequest]{
// 		Msg: &apiv1.IPServiceUpdateRequest{
// 			Project: rq.Project,
// 			Ip:      rq,
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Msg.Ip, nil
// }

func ipResponseToCreate(r *apiv1.IP) *connect.Request[apiv1.IPServiceAllocateRequest] {
	return &connect.Request[apiv1.IPServiceAllocateRequest]{
		Msg: &apiv1.IPServiceAllocateRequest{
			Project:     r.Project,
			Name:        r.Name,
			Description: r.Description,
			Tags:        r.Tags,
			Static:      ipTypeToStatic(r.Type),
		},
	}
}

func ipResponseToUpdate(r *apiv1.IP) *apiv1.IP {
	return r
}

func ipStaticToType(b bool) apiv1.IPType {
	if b {
		return apiv1.IPType_IP_TYPE_STATIC
	}
	return apiv1.IPType_IP_TYPE_EPHEMERAL
}

func ipTypeToStatic(t apiv1.IPType) bool {
	return t == apiv1.IPType_IP_TYPE_STATIC
}
