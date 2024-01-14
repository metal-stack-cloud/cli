package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func ClusterSorter() *multisort.Sorter[*apiv1.Cluster] {
	return multisort.New(multisort.FieldMap[*apiv1.Cluster]{
		"uptime": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CreatedAt.AsTime().UnixMilli(), b.CreatedAt.AsTime().UnixMilli(), descending)
		},
		"name": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"project": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"uuid": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"partition": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Partition, b.Partition, descending)
		},
		"tenant": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Tenant, b.Tenant, descending)
		},
		"purpose": func(a, b *apiv1.Cluster, descending bool) multisort.CompareResult {
			return multisort.Compare(pointer.SafeDeref(a.Purpose), pointer.SafeDeref(b.Purpose), descending)
		},
	}, multisort.Keys{{ID: "tenant"}, {ID: "project"}, {ID: "name"}, {ID: "uuid"}})
}
