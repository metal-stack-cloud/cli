package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func ClusterSorter() *multisort.Sorter[*apiv1.Cluster] {
	return multisort.New(multisort.FieldMap[*apiv1.Cluster]{
		"uptime": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CreatedAt.AsTime().UnixMilli(), b.CreatedAt.AsTime().UnixMilli(), descending)
		},
		"name": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
	}, multisort.Keys{{ID: "uptime", Descending: true}})
}
