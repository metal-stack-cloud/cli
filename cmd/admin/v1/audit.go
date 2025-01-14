package v1

import (
	"fmt"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
)

type audit struct {
	c *config.Config
}

func newAuditCmd(c *config.Config) *cobra.Command {
	a := &audit{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.AuditServiceGetRequest, *apiv1.AuditServiceListRequest, *apiv1.AuditTrace]{
		BinaryName:      config.BinaryName,
		Singular:        "audit trace",
		Plural:          "audit traces",
		Description:     "show audit traces of the api-server",
		Sorter:          sorters.AuditSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list all audit traces",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "gets the audit trace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return genericcli.NewCmds(cmdsConfig, listCmd, getCmd)
}

func (a *audit) Get(req *apiv1.AuditServiceGetRequest) (*apiv1.AuditTrace, error) {
	ctx, cancel := a.c.NewRequestContext()
	defer cancel()

	// not sure about how to get tenant of current user
	tenant, err := a.c.GetTenant()
	if err != nil {
		return nil, fmt.Errorf("tenant is required")
	}
	req.Login = tenant

	resp, err := a.c.Client.Apiv1().Audit().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get audit trace: %w", err)
	}

	return resp.Msg.Audit, nil
}

func (a *audit) List(req *apiv1.AuditServiceListRequest) ([]*apiv1.AuditTrace, error) {
	ctx, cancel := a.c.NewRequestContext()
	defer cancel()

	resp, err := a.c.Client.Apiv1().Audit().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list audit traces: %w", err)
	}

	return resp.Msg.Audits, nil
}
