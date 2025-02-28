package v1

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type audit struct {
	c *config.Config
}

func newAuditCmd(c *config.Config) *cobra.Command {
	a := &audit{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.AuditServiceGetRequest, *apiv1.AuditServiceGetRequest, *apiv1.AuditTrace]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*apiv1.AuditServiceGetRequest, *apiv1.AuditServiceGetRequest, *apiv1.AuditTrace](a).WithFS(c.Fs),
		Singular:        "audit trace",
		Plural:          "audit traces",
		Description:     "show audit traces of the api-server",
		Sorter:          sorters.AuditSorter(),
		OnlyCmds:        genericcli.OnlyCmds(genericcli.ListCmd, genericcli.DescribeCmd),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {

			cmd.Flags().String("request-id", "", "request id of the audit trace.")

			cmd.Flags().String("from", "1h", "start of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05")
			cmd.Flags().String("to", "", "end of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05")

			cmd.Flags().String("user", "", "user of the audit trace.")
			cmd.Flags().String("tenant", "", "tenant of the audit trace.")

			cmd.Flags().String("project", "", "project id of the audit trace")

			cmd.Flags().String("method", "", "api method of the audit trace.")
			cmd.Flags().Int32("result-code", 0, "HTTP status code of the audit trace.")
			cmd.Flags().String("source-ip", "", "source-ip of the audit trace.")

			cmd.Flags().String("body", "", "filters audit trace body payloads for the giben text.")
			cmd.Flags().String("error", "", "error of the audit trace.")

			//removed since issues arise with current flow of merging req and res to one request
			//cmd.Flags().Int64("limit", 100, "limit the number of audit traces.")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().Bool("prettify-body", false, "attempts to interpret the body as json and prettifies it.")
		},
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (a *audit) Get(id string) (*apiv1.AuditTrace, error) {
	ctx, cancel := a.c.NewRequestContext()
	defer cancel()

	// not sure about how to get tenant of current user
	tenant, err := a.c.GetTenant()
	if err != nil {
		return nil, fmt.Errorf("tenant is required")
	}

	req := &apiv1.AuditServiceGetRequest{
		Login: tenant,
		Uuid:  id,
	}

	resp, err := a.c.Client.Apiv1().Audit().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get audit trace: %w", err)
	}

	trace := resp.Msg.Audit

	if viper.GetBool("prettify-body") {
		trimmed := strings.Trim(trace.RequestPayload, `"`)
		body := map[string]any{}
		err = json.Unmarshal([]byte(trimmed), &body)
		if err == nil {
			if pretty, err := json.MarshalIndent(body, "", "    "); err == nil {
				trace.RequestPayload = string(pretty)
			}
		}
		trimmed = strings.Trim(trace.ResponsePayload, `"`)
		err = json.Unmarshal([]byte(trimmed), &body)
		if err == nil {
			if pretty, err := json.MarshalIndent(body, "", "    "); err == nil {
				trace.ResponsePayload = string(pretty)
			}
		}
	}

	return resp.Msg.Audit, nil
}

func (a *audit) List() ([]*apiv1.AuditTrace, error) {
	ctx, cancel := a.c.NewRequestContext()
	defer cancel()

	fromDateTime, err := eventuallyRelativeDateTime(viper.GetString("from"))
	if err != nil {
		return nil, err
	}
	toDateTime, err := eventuallyRelativeDateTime(viper.GetString("to"))
	if err != nil {
		return nil, err
	}

	tenant, err := a.c.GetTenant()
	if err != nil {
		return nil, fmt.Errorf("tenant is required %w", err)
	}

	req := &apiv1.AuditServiceListRequest{
		Login:      tenant,
		Uuid:       pointer.PointerOrNil(viper.GetString("request-id")),
		From:       fromDateTime,
		To:         toDateTime,
		User:       pointer.PointerOrNil(viper.GetString("user")),
		Project:    pointer.PointerOrNil(viper.GetString("project")),
		Method:     pointer.PointerOrNil(viper.GetString("method")),
		ResultCode: pointer.PointerOrNil(viper.GetInt32("result-code")),
		Body:       pointer.PointerOrNil(viper.GetString("body")),
		SourceIp:   pointer.PointerOrNil(viper.GetString("source-ip")),
	}

	resp, err := a.c.Client.Apiv1().Audit().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to list audit traces: %w", err)
	}

	return resp.Msg.Audits, nil
}

func eventuallyRelativeDateTime(s string) (*timestamppb.Timestamp, error) {
	if s == "" {
		return timestamppb.Now(), nil
	}
	duration, err := time.ParseDuration(s)
	if err == nil {
		return timestamppb.New(time.Now().Add(-duration)), nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return timestamppb.Now(), fmt.Errorf("failed to convert time: %w", err)
	}
	return timestamppb.New(t), nil
}

func (a *audit) Convert(*apiv1.AuditTrace) (string, *apiv1.AuditServiceGetRequest, *apiv1.AuditServiceGetRequest, error) {
	return "", nil, nil, fmt.Errorf("not implemented for audit traces")
}

func (a *audit) Delete(id string) (*apiv1.AuditTrace, error) {
	return nil, fmt.Errorf("not implemented for audit traces")
}

func (a *audit) Create(*apiv1.AuditServiceGetRequest) (*apiv1.AuditTrace, error) {
	return nil, fmt.Errorf("not implemented for audit traces")
}

func (a *audit) Update(*apiv1.AuditServiceGetRequest) (*apiv1.AuditTrace, error) {
	return nil, fmt.Errorf("not implemented for audit traces")
}
