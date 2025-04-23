package tableprinters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/grpc/codes"
)

func (t *TablePrinter) AuditTable(data []*apiv1.AuditTrace, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"Time", "Request-Id", "User", "Project", "Method", "Phase", "Code"}
	if wide {
		header = []string{"Time", "Request-Id", "User", "Project", "Method", "Phase", "Source-Ip", "Code", "Body"}
	}

	for _, audit := range data {
		id := audit.Uuid
		time := audit.Timestamp.AsTime().Format("2006-01-02 15:04:05")
		user := audit.User
		phase := audit.Phase
		method := audit.Method
		sourceIp := audit.SourceIp
		project := pointer.SafeDeref(audit.Project)
		body := genericcli.TruncateEnd(pointer.SafeDeref(audit.Body), 30)

		code := ""
		if audit.ResultCode != nil {
			code = codes.Code(uint32(*audit.ResultCode)).String()
		}

		if wide {
			rows = append(rows, []string{time, id, user, project, method, phase.String(), sourceIp, code, body})
		} else {
			rows = append(rows, []string{time, id, user, project, method, phase.String(), code})
		}
	}

	t.t.MutateTable(func(table *tablewriter.Table) {
		table.SetAutoWrapText(false)
	})

	return header, rows, nil
}
