package v1

import (
	"fmt"

	"connectrpc.com/connect"
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
			cmd.Flags().StringP("project", "p", "", "filter by project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionAssetListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "filter by project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "filter by project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		OnlyCmds: genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DeleteCmd, genericcli.DescribeCmd),
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (s *snapshot) Create(rq any) (*apiv1.Snapshot, error) {
	panic("unimplemented")
}

func (c *snapshot) Delete(id string) (*apiv1.Snapshot, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.SnapshotServiceDeleteRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Snapshot().Delete(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to delete snapshots: %w", err)
	}

	return resp.Msg.Snapshot, nil
}

func (c *snapshot) Get(id string) (*apiv1.Snapshot, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.SnapshotServiceGetRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Snapshot().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	return resp.Msg.Snapshot, nil
}

func (c *snapshot) List() ([]*apiv1.Snapshot, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.SnapshotServiceListRequest{
		Project: c.c.GetProject(),
	}
	if viper.IsSet("uuid") {
		req.Uuid = pointer.Pointer(viper.GetString("uuid"))
	}
	if viper.IsSet("name") {
		req.Name = pointer.Pointer(viper.GetString("name"))
	}
	if viper.IsSet("partition") {
		req.Partition = pointer.Pointer(viper.GetString("partition"))
	}

	resp, err := c.c.Client.Apiv1().Snapshot().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshots: %w", err)
	}

	return resp.Msg.Snapshots, nil
}

func (c *snapshot) Convert(r *apiv1.Snapshot) (string, any, any, error) {
	panic("unimplemented")
}

func (c *snapshot) Update(rq any) (*apiv1.Snapshot, error) {
	panic("unimplemented")
}
