package v1

import (
	"fmt"

	"connectrpc.com/connect"
	adminApiv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	helpersaudit "github.com/metal-stack-cloud/cli/pkg/helpers/audit"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type adminAudit struct {
	c *config.Config
}

func newAuditCmd(c *config.Config) *cobra.Command {
	a := &adminAudit{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.AuditTrace]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI(a).WithFS(c.Fs),
		Singular:        "audit",
		Plural:          "audits",
		Description:     "show audit traces of the api-server",
		Sorter:          sorters.AuditSorter(),
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("request-id", "", "request id of the audit trace.")

			cmd.Flags().String("from", "", "start of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05")
			cmd.Flags().String("to", "", "end of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05")

			cmd.Flags().String("user", "", "user of the audit trace.")
			cmd.Flags().String("tenant", "", "tenant of the audit trace.")

			cmd.Flags().String("project", "", "project id of the audit trace")

			cmd.Flags().String("phase", "", "the audit trace phase.")
			cmd.Flags().String("method", "", "api method of the audit trace.")
			cmd.Flags().Int32("result-code", 0, "gRPC result status code of the audit trace.")
			cmd.Flags().String("source-ip", "", "source-ip of the audit trace.")

			cmd.Flags().String("body", "", "filters audit trace body payloads for the given text (full-text search).")

			cmd.Flags().Int64("limit", 0, "limit the number of audit traces.")

			cmd.Flags().Bool("prettify-body", false, "attempts to interpret the body as json and prettifies it.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("phase", c.Completion.AuditPhaseListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("result-code", c.Completion.AuditStatusCodesCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("tenant", "", "tenant of the audit trace.")

			cmd.Flags().String("phase", "", "the audit trace phase.")

			cmd.Flags().Bool("prettify-body", false, "attempts to interpret the body as json and prettifies it.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("phase", c.Completion.AuditPhaseListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("tenant", c.Completion.TenantListCompletion))
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (a *adminAudit) Get(id string) (*apiv1.AuditTrace, error) {
	// NOTE: required by api but not needed here since admin audit is only used for listing audits
	panic("unimplemented")
}

func (a *adminAudit) List() ([]*apiv1.AuditTrace, error) {
	ctx, cancel := a.c.NewRequestContext()
	defer cancel()

	fromDateTime, err := helpersaudit.EventuallyRelativeDateTime(viper.GetString("from"))
	if err != nil {
		return nil, err
	}
	toDateTime, err := helpersaudit.EventuallyRelativeDateTime(viper.GetString("to"))
	if err != nil {
		return nil, err
	}

	var code *int32
	if viper.IsSet("result-code") {
		code = pointer.Pointer(viper.GetInt32("result-code"))
	}

	req := &adminApiv1.AuditServiceListRequest{
		Uuid:       pointer.PointerOrNil(viper.GetString("request-id")),
		From:       fromDateTime,
		To:         toDateTime,
		User:       pointer.PointerOrNil(viper.GetString("user")),
		Project:    pointer.PointerOrNil(viper.GetString("project")),
		Method:     pointer.PointerOrNil(viper.GetString("method")),
		ResultCode: code,
		Body:       pointer.PointerOrNil(viper.GetString("body")),
		SourceIp:   pointer.PointerOrNil(viper.GetString("source-ip")),
		Limit:      pointer.PointerOrNil(viper.GetInt32("limit")),
		Phase:      a.toPhase(viper.GetString("phase")),
	}

	resp, err := a.c.Client.Adminv1().Audit().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list audit traces: %w", err)
	}

	if viper.GetBool("prettify-body") {
		for _, trace := range resp.Msg.Traces {
			a.tryPrettifyBody(trace)
		}
	}

	return resp.Msg.Traces, nil
}

func (a *adminAudit) Convert(*apiv1.AuditTrace) (string, any, any, error) {
	return helpersaudit.Convert()
}

func (a *adminAudit) Delete(id string) (*apiv1.AuditTrace, error) {
	return helpersaudit.Delete()
}

func (a *adminAudit) Create(any) (*apiv1.AuditTrace, error) {
	return helpersaudit.Create()
}

func (a *adminAudit) Update(any) (*apiv1.AuditTrace, error) {
	return helpersaudit.Update()
}

func (a *adminAudit) toPhase(phase string) *apiv1.AuditPhase {
	return helpersaudit.ToPhase(phase)
}

func (a *adminAudit) tryPrettifyBody(trace *apiv1.AuditTrace) {
	helpersaudit.TryPrettifyBody(trace)
}
