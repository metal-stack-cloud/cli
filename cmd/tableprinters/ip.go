package tableprinters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) IPTable(data []*apiv1.IP, _ bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"IP", "ID", "Project", "Name", "Description", "Type"}
	)

	for _, ip := range data {
		var t string
		switch ip.Type {
		case apiv1.IPType_IP_TYPE_EPHEMERAL:
			t = "ephemeral"
		case apiv1.IPType_IP_TYPE_STATIC:
			t = "static"
		case apiv1.IPType_IP_TYPE_UNSPECIFIED:
			t = "unspecified"
		default:
			t = ip.Type.String()
		}

		rows = append(rows, []string{ip.Ip, ip.Uuid, ip.Project, ip.Name, ip.Description, t})
	}

	return header, rows, nil
}
