package tableprinters

import (
	"strings"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/tag"
)

func (t *TablePrinter) IPTable(data []*apiv1.IP, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"IP", "Project", "ID", "Type", "Name", "Attached Service"}
	)

	if wide {
		header = []string{"IP", "Project", "ID", "Type", "Name", "Description", "Labels"}
	}

	for _, ip := range data {
		ip := ip

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

		attachedService := ""

		tm := tag.NewTagMap(ip.Tags)
		if value, ok := tm.Value(tag.ClusterServiceFQN); ok {
			if parts := strings.Split(value, "/"); len(parts) == 3 {
				attachedService = parts[2]
			}
		}

		if wide {
			rows = append(rows, []string{ip.Ip, ip.Project, ip.Uuid, t, ip.Name, ip.Description, strings.Join(ip.Tags, "\n")})
		} else {
			rows = append(rows, []string{ip.Ip, ip.Project, ip.Uuid, t, ip.Name, attachedService})
		}
	}
	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
