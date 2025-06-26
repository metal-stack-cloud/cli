package tableprinters

import (
	"fmt"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

type TablePrinter struct {
	t *printers.TablePrinter
}

func New() *TablePrinter {
	return &TablePrinter{}
}

func (t *TablePrinter) SetPrinter(printer *printers.TablePrinter) {
	t.t = printer
}

func (t *TablePrinter) ToHeaderAndRows(data any, wide bool) ([]string, [][]string, error) {
	switch d := data.(type) {

	case *apiv1.Asset:
		return t.AssetTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Asset:
		return t.AssetTable(d, wide)

	case *apiv1.AuditTrace:
		return t.AuditTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.AuditTrace:
		return t.AuditTable(d, wide)

	case *config.Contexts:
		return t.ContextTable(d, wide)

	case *apiv1.IP:
		return t.IPTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.IP:
		return t.IPTable(d, wide)

	case *apiv1.Coupon:
		return t.CouponTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Coupon:
		return t.CouponTable(d, wide)

	case *apiv1.Cluster:
		return t.ClusterTable(pointer.WrapInSlice(d), nil, wide)
	case []*apiv1.Cluster:
		return t.ClusterTable(d, nil, wide)
	case *adminv1.ClusterServiceGetResponse:
		return t.ClusterMachineTable(pointer.WrapInSlice(d), wide)
	case *apiv1.ClusterStatusLastError:
		return t.ClusterStatusLastErrorTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.ClusterStatusLastError:
		return t.ClusterStatusLastErrorTable(d, wide)
	case *apiv1.ClusterStatusCondition:
		return t.ClusterStatusConditionsTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.ClusterStatusCondition:
		return t.ClusterStatusConditionsTable(d, wide)

	case *apiv1.Price:
		return t.PaymentPricesTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Price:
		return t.PaymentPricesTable(d, wide)

	case *apiv1.Project:
		return t.ProjectTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Project:
		return t.ProjectTable(d, wide)
	case *apiv1.ProjectInvite:
		return t.ProjectInviteTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.ProjectInvite:
		return t.ProjectInviteTable(d, wide)
	case *apiv1.ProjectMember:
		return t.ProjectMemberTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.ProjectMember:
		return t.ProjectMemberTable(d, wide)

	case *apiv1.Volume:
		return t.VolumeTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Volume:
		return t.VolumeTable(d, wide)

	case *apiv1.Snapshot:
		return t.SnapshotTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Snapshot:
		return t.SnapshotTable(d, wide)

	case *adminv1.StorageClusterInfo:
		return t.StorageClusterInfoTable(pointer.WrapInSlice(d), wide)
	case []*adminv1.StorageClusterInfo:
		return t.StorageClusterInfoTable(d, wide)

	case *apiv1.Token:
		return t.TokenTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Token:
		return t.TokenTable(d, wide)

	case *apiv1.Tenant:
		return t.TenantTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Tenant:
		return t.TenantTable(d, wide)
	case *apiv1.TenantInvite:
		return t.TenantInviteTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.TenantInvite:
		return t.TenantInviteTable(d, wide)
	case *apiv1.TenantMember:
		return t.TenantMemberTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.TenantMember:
		return t.TenantMemberTable(d, wide)

	case *apiv1.Health:
		return t.HealthTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Health:
		return t.HealthTable(d, wide)

	default:
		return nil, nil, fmt.Errorf("unknown table printer for type: %T", d)
	}
}
