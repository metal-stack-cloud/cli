package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func TenantSorter() *multisort.Sorter[*apiv1.Tenant] {
	return multisort.New(multisort.FieldMap[*apiv1.Tenant]{
		"id": func(a, b *apiv1.Tenant, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Login, b.Login, descending)
		},
		"name": func(a, b *apiv1.Tenant, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"since": func(a, b *apiv1.Tenant, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CreatedAt.AsTime().UnixMilli(), b.CreatedAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "since", Descending: true}})
}
