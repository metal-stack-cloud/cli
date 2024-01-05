package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func ProjectSorter() *multisort.Sorter[*apiv1.Project] {
	return multisort.New(multisort.FieldMap[*apiv1.Project]{
		"id": func(a, b *apiv1.Project, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"name": func(a, b *apiv1.Project, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"tenant": func(a, b *apiv1.Project, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Tenant, b.Tenant, descending)
		},
	}, multisort.Keys{{ID: "tenant"}, {ID: "name"}, {ID: "id"}})
}
