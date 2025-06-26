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

func TenantInviteSorter() *multisort.Sorter[*apiv1.TenantInvite] {
	return multisort.New(multisort.FieldMap[*apiv1.TenantInvite]{
		"tenant": func(a, b *apiv1.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Tenant, b.Tenant, descending)
		},
		"secret": func(a, b *apiv1.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Secret, b.Secret, descending)
		},
		"role": func(a, b *apiv1.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"expiration": func(a, b *apiv1.TenantInvite, descending bool) multisort.CompareResult {
			return multisort.Compare(a.ExpiresAt.AsTime().UnixMilli(), b.ExpiresAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "tenant"}, {ID: "role"}, {ID: "expiration"}})
}

func TenantMemberSorter() *multisort.Sorter[*apiv1.TenantMember] {
	return multisort.New(multisort.FieldMap[*apiv1.TenantMember]{
		"id": func(a, b *apiv1.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Id, b.Id, descending)
		},
		"role": func(a, b *apiv1.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.Role, b.Role, descending)
		},
		"created": func(a, b *apiv1.TenantMember, descending bool) multisort.CompareResult {
			return multisort.Compare(a.CreatedAt.AsTime().UnixMilli(), b.CreatedAt.AsTime().UnixMilli(), descending)
		},
	}, multisort.Keys{{ID: "role"}, {ID: "id"}})
}
