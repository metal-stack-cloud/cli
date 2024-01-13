package v1

import (
	"fmt"

	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
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
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "filter by uuid")
			cmd.Flags().StringP("name", "", "", "filter by name")
			cmd.Flags().StringP("partition", "", "", "filter by partition")
			cmd.Flags().StringP("project", "p", "", "filter by project")
			cmd.Flags().StringP("tenant", "", "", "filter by tenant")
		},
		OnlyCmds: genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd),
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (v *volume) Create(rq any) (*apiv1.Volume, error) {
	panic("unimplemented")
}

func (v *volume) Delete(id string) (*apiv1.Volume, error) {
	panic("unimplemented")
}

func (c *volume) Get(id string) (*apiv1.Volume, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv1.StorageServiceListVolumesRequest{
		Uuid: &id,
	}

	resp, err := c.c.Client.Adminv1().Storage().ListVolumes(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}

	if len(resp.Msg.Volumes) != 1 {
		return nil, fmt.Errorf("no volume with ID:%s found", id)
	}

	return resp.Msg.Volumes[0], nil
}

func (c *volume) List() ([]*apiv1.Volume, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv1.StorageServiceListVolumesRequest{}

	if viper.IsSet("uuid") {
		req.Uuid = pointer.Pointer(viper.GetString("uuid"))
	}
	if viper.IsSet("name") {
		req.Name = pointer.Pointer(viper.GetString("name"))
	}
	if viper.IsSet("project") {
		req.Project = pointer.Pointer(viper.GetString("project"))
	}
	if viper.IsSet("partition") {
		req.Partition = pointer.Pointer(viper.GetString("partition"))
	}
	if viper.IsSet("tenant") {
		req.Tenant = pointer.Pointer(viper.GetString("tenant"))
	}

	resp, err := c.c.Client.Adminv1().Storage().ListVolumes(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}

	return resp.Msg.Volumes, nil
}

func (v *volume) Convert(r *apiv1.Volume) (string, any, any, error) {
	panic("unimplemented")
}

func (v *volume) Update(rq any) (*apiv1.Volume, error) {
	panic("unimplemented")
}
