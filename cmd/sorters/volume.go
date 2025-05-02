package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func VolumeSorter() *multisort.Sorter[*apiv1.Volume] {
	return multisort.New(multisort.FieldMap[*apiv1.Volume]{
		"uuid": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"name": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"project": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"partition": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Partition, b.Partition, descending)
		},
		"storage-class": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.StorageClass, b.StorageClass, descending)
		},
		"size": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Size, b.Size, descending)
		},
		"usage": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Usage, b.Usage, descending)
		},
		"state": func(a, b *apiv1.Volume, descending bool) multisort.CompareResult {
			return multisort.Compare(a.State, b.State, descending)
		},
	}, multisort.Keys{{ID: "name"}, {ID: "project"}, {ID: "size"}, {ID: "storage-class"}, {ID: "partition"}, {ID: "usage"}, {ID: "state"}, {ID: "uuid"}})
}
