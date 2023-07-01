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

type snapshot struct {
	c *config.Config
}

func newSnapshotCmd(c *config.Config) *cobra.Command {
	w := &snapshot{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Snapshot]{
		BinaryName:  config.BinaryName,
		GenericCLI:  genericcli.NewGenericCLI[any, any, *apiv1.Snapshot](w).WithFS(c.Fs),
		Singular:    "snapshot",
		Plural:      "snapshots",
		Description: "snapshot related actions of metalstack.cloud",
		// Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
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
		OnlyCmds: genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DeleteCmd, genericcli.DescribeCmd),
	}
	return genericcli.NewCmds(cmdsConfig)
}

// Create implements genericcli.CRUD
func (s *snapshot) Create(rq any) (*apiv1.Snapshot, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (s *snapshot) Delete(id string) (*apiv1.Snapshot, error) {
	req := &apiv1.SnapshotServiceDeleteRequest{
		Uuid:    id,
		Project: viper.GetString("project"),
	}
	resp, err := s.c.Client.Apiv1().Snapshot().Delete(s.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to delete snapshots: %w", err)
	}
	return resp.Msg.Snapshot, nil
}

// Get implements genericcli.CRUD
func (s *snapshot) Get(id string) (*apiv1.Snapshot, error) {
	req := &apiv1.SnapshotServiceGetRequest{
		Uuid:    id,
		Project: viper.GetString("project"),
	}
	resp, err := s.c.Client.Apiv1().Snapshot().Get(s.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}
	return resp.Msg.Snapshot, nil
}

// List implements genericcli.CRUD
func (s *snapshot) List() ([]*apiv1.Snapshot, error) {
	req := &apiv1.SnapshotServiceListRequest{}
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
	resp, err := s.c.Client.Apiv1().Snapshot().List(s.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}
	return resp.Msg.Snapshots, nil
}

// Convert implements genericcli.CRUD
func (s *snapshot) Convert(r *apiv1.Snapshot) (string, any, any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (s *snapshot) Update(rq any) (*apiv1.Snapshot, error) {
	panic("unimplemented")
}
