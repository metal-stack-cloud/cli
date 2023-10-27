package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func TokenSorter() *multisort.Sorter[*apiv1.Token] {
	return multisort.New(multisort.FieldMap[*apiv1.Token]{
		"id": func(a, b *apiv1.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
		"expires": func(a, b *apiv1.Token, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Expires.AsTime().UnixMilli(), b.Expires.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "id"}})
}
