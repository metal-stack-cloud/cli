package tableprinters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) IPTable(data []*apiv1.IP, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"IP", "Project"}
	if wide {
		header = []string{"IP", "Project"}
	}

	for _, ip := range data {
		if wide {
			rows = append(rows, []string{ip.Ip, ip.Project})
		} else {
			rows = append(rows, []string{ip.Ip, ip.Project})
		}
	}

	return header, rows, nil
}
