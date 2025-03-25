package tableprinters

import (
	"fmt"
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
)

func (t *TablePrinter) AuditTable(data []*apiv1.AuditTrace, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"TIME", "REQUEST-ID", "USER", "PROJECT", "METHOD", "PHASE"}
	if wide {
		header = []string{"TIME", "REQUEST-ID", "USER", "PROJECT", "METHOD", "PHASE", "SOURCE-IP", "RESULT-CODE", "BODY"}
	}

	for _, audit := range data {
		id := audit.Uuid
		time := truncateToSeconds(audit.Timestamp.AsTime()).Format("2006-01-02 15:04:05")
		user := audit.User
		phase := audit.Phase
		method := audit.Method
		sourceIp := audit.SourceIp

		project := ""
		if audit.Project != nil {
			project = *audit.Project
		}
		body := ""
		if audit.Body != nil {
			body = genericcli.TruncateEnd(*audit.Body, 30)
		}

		if wide {
			var resultCode = ""
			if audit.ResultCode != nil {
				resultCode = fmt.Sprintf("%d", audit.ResultCode)
			}
			rows = append(rows, []string{time, id, user, project, method, phase, sourceIp, resultCode, body})
		} else {
			rows = append(rows, []string{time, id, user, project, method, phase})
		}
	}

	return header, rows, nil
}

func truncateToSeconds(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}
