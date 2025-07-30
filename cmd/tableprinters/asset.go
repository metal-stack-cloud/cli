package tableprinters

import (
	"strings"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) AssetTable(data *apiv1.AssetServiceListResponse, _ bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Region", "Partition", "Machine Types"}
	)

	for _, asset := range data.Assets {
		if asset.Region == nil {
			continue
		}

		region := asset.Region

		if !region.Active {
			continue
		}

		var machineTypes []string
		for _, p := range asset.MachineTypes {
			p := p

			machineTypes = append(machineTypes, p.Id)
		}

		for _, p := range region.Partitions {
			p := p

			rows = append(rows, []string{region.Id, p.Id, strings.Join(machineTypes, ",")})
		}

	}

	return header, rows, nil
}
