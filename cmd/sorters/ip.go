package sorters

import (
	"net/netip"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func IPSorter() *multisort.Sorter[*apiv1.IP] {
	return multisort.New(multisort.FieldMap[*apiv1.IP]{
		"ip": func(a, b *apiv1.IP, descending bool) multisort.CompareResult {
			aIP, _ := netip.ParseAddr(a.Ip)
			bIP, _ := netip.ParseAddr(b.Ip)
			return multisort.WithCompareFunc(func() int {
				return aIP.Compare(bIP)
			}, descending)
		},
		"name": func(a, b *apiv1.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Name, b.Name, descending)
		},
		"project": func(a, b *apiv1.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"type": func(a, b *apiv1.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Type, b.Type, descending)
		},
		"network": func(a, b *apiv1.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Network, b.Network, descending)
		},
		"uuid": func(a, b *apiv1.IP, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Uuid, b.Uuid, descending)
		},
	}, multisort.Keys{{ID: "project"}, {ID: "ip"}})
}
