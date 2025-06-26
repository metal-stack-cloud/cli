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

func ProjectInviteSorter() *multisort.Sorter[*apiv1.ProjectInvite] {
	return multisort.New(multisort.FieldMap[*apiv1.ProjectInvite]{
		"project": func(a, b *apiv1.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Project, b.Project, descending)
		},
		"secret": func(a, b *apiv1.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Secret, b.Secret, descending)
		},
		"role": func(a, b *apiv1.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"expiration": func(a, b *apiv1.ProjectInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ExpiresAt.AsTime().UnixMilli(), b.ExpiresAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "project"}, {ID: "role"}, {ID: "expiration"}})
}

func ProjectMemberSorter() *multisort.Sorter[*apiv1.ProjectMember] {
	return multisort.New(multisort.FieldMap[*apiv1.ProjectMember]{
		"id": func(a, b *apiv1.ProjectMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"role": func(a, b *apiv1.ProjectMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"created": func(a, b *apiv1.ProjectMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CreatedAt.AsTime().UnixMilli(), b.CreatedAt.AsTime().UnixMilli(), descending)
		},
		"inherited": func(a, b *apiv1.ProjectMember, descending bool) multisort.CompareResult {
			boolToInt := func(in bool) int {
				if in {
					return 1
				}
				return 0
			}
			return multisort.Compare(boolToInt(a.InheritedMembership), boolToInt(b.InheritedMembership), descending)
		},
	}, multisort.Keys{{ID: "inherited", Descending: false}, {ID: "role"}, {ID: "id"}})
}
