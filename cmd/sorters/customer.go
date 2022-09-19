package sorters

import (
	v1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
	p "github.com/metal-stack/metal-lib/pkg/pointer"
)

func CustomerSorter() *multisort.Sorter[*v1.PaymentCustomer] {
	return multisort.New(multisort.FieldMap[*v1.PaymentCustomer]{
		"id": func(a, b *v1.PaymentCustomer, descending bool) multisort.CompareResult {
			return multisort.Compare(p.SafeDeref(a.CustomerId), p.SafeDeref(b.CustomerId), descending)
		},
		"name": func(a, b *v1.PaymentCustomer, descending bool) multisort.CompareResult {
			return multisort.Compare(p.SafeDeref(a.Name), p.SafeDeref(b.Name), descending)
		},
	}, multisort.Keys{{ID: "name"}, {ID: "id"}})
}
