package sorters

import (
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func UserSorter() *multisort.Sorter[*adminv1.User] {
	return multisort.New(multisort.FieldMap[*adminv1.User]{
		"id": func(a, b *adminv1.User, descending bool) multisort.CompareResult {
			return multisort.Compare(a.User.Login, b.User.Login, descending)
		},
		"name": func(a, b *adminv1.User, descending bool) multisort.CompareResult {
			return multisort.Compare(a.User.Name, b.User.Name, descending)
		},
	}, multisort.Keys{{ID: "name"}, {ID: "id"}})
}
