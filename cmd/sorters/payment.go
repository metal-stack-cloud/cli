package sorters

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func PriceSorter() *multisort.Sorter[*apiv1.Price] {
	return multisort.New(multisort.FieldMap[*apiv1.Price]{
		"type": func(a, b *apiv1.Price, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ProductType.String(), b.ProductType.String(), descending)
		},
		"name": func(a, b *apiv1.Price, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"price": func(a, b *apiv1.Price, descending bool) multisort.CompareResult {
			return multisort.Compare(a.GetUnitAmountDecimal(), b.GetUnitAmountDecimal(), descending)
		},
	}, multisort.Keys{{ID: "type"}, {ID: "name"}, {ID: "price"}})
}
