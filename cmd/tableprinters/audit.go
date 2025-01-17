package tableprinters

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) AuditTable(data []*apiv1.AuditTrace, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"TIME", "REQUEST-ID", "USER", "METHOD"}
	if wide {
		header = []string{"TIME", "REQUEST-ID", "USER", "PROJECT", "METHOD", "SOURCE-IP", "RESULT-CODE", "BODY"}
	}

	for _, audit := range data {
		id := audit.Uuid
		time := truncateToSeconds(audit.Timestamp.AsTime()).Format("2006-01-02 15:04:05")
		user := audit.User
		project := audit.Project

		method := audit.Method
		sourceIp := audit.SourceIp
		resultCode := audit.ResultCode

		body := audit.ResponsePayload

		if wide {
			rows = append(rows, []string{time, id, user, project, method, sourceIp, resultCode, body})
		} else {
			rows = append(rows, []string{time, id, user, method})
		}
	}

	return header, rows, nil
}

func truncateToSeconds(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}
