package sorters

import (
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func ContextSorter() *multisort.Sorter[*config.Context] {
	return multisort.New(multisort.FieldMap[*config.Context]{
		"name": func(a, b *config.Context, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
	}, multisort.Keys{{ID: "name"}})
}
