package sorters

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/multisort"
)

func AuditSorter() *multisort.Sorter[*apiv1.AuditTrace] {
	return multisort.New(multisort.FieldMap[*apiv1.AuditTrace]{
		"timestamp": func(a, b *apiv1.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(time.Time(a.Timestamp.AsTime()).Unix(), time.Time(b.Timestamp.AsTime()).Unix(), descending)
		},
		"user": func(a, b *apiv1.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(a.User, b.User, descending)
		},
		"method": func(a, b *apiv1.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Method, b.Method, descending)
		},
		"project": func(a, b *apiv1.AuditTrace, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
	}, multisort.Keys{{ID: "timestamp", Descending: true}})
}
