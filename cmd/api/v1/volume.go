package v1

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type volume struct {
	c *config.Config
}

func newVolumeCmd(c *config.Config) *cobra.Command {
	w := &volume{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Volume]{
		BinaryName:  config.BinaryName,
		GenericCLI:  genericcli.NewGenericCLI[any, any, *apiv1.Volume](w).WithFS(c.Fs),
		Singular:    "volume",
		Plural:      "volumes",
		Description: "volume related actions of metalstack.cloud",
		// Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "filter by uuid")
			cmd.Flags().StringP("name", "", "", "filter by name")
			cmd.Flags().StringP("partition", "", "", "filter by partition")
			cmd.Flags().StringP("project", "", "", "filter by project")
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "filter by uuid")
			cmd.Flags().StringP("project", "", "", "filter by project")
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "filter by uuid")
			cmd.Flags().StringP("project", "", "", "filter by project")
		},
	}
	return genericcli.NewCmds(cmdsConfig)
}

// Create implements genericcli.CRUD
func (v *volume) Create(rq any) (*apiv1.Volume, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (v *volume) Delete(id string) (*apiv1.Volume, error) {
	req := &apiv1.VolumeServiceDeleteRequest{
		Uuid:    id,
		Project: viper.GetString("project"),
	}
	resp, err := v.c.Apiv1Client.Volume().Delete(v.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to delete volumes: %w", err)
	}
	return resp.Msg.Volume, nil
}

// Get implements genericcli.CRUD
func (v *volume) Get(id string) (*apiv1.Volume, error) {
	req := &apiv1.VolumeServiceGetRequest{
		Uuid:    id,
		Project: viper.GetString("project"),
	}
	resp, err := v.c.Apiv1Client.Volume().Get(v.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}
	return resp.Msg.Volume, nil
}

// List implements genericcli.CRUD
func (v *volume) List() ([]*apiv1.Volume, error) {
	req := &apiv1.VolumeServiceListRequest{}
	if viper.IsSet("uuid") {
		req.Uuid = pointer.Pointer(viper.GetString("uuid"))
	}
	if viper.IsSet("name") {
		req.Name = pointer.Pointer(viper.GetString("name"))
	}
	if viper.IsSet("project") {
		req.Project = viper.GetString("project")
	}
	if viper.IsSet("partition") {
		req.Partition = pointer.Pointer(viper.GetString("partition"))
	}
	resp, err := v.c.Apiv1Client.Volume().List(v.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}
	return resp.Msg.Volumes, nil
}

// ToCreate implements genericcli.CRUD
func (v *volume) ToCreate(r *apiv1.Volume) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (v *volume) ToUpdate(r *apiv1.Volume) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (v *volume) Update(rq any) (*apiv1.Volume, error) {
	panic("unimplemented")
}
