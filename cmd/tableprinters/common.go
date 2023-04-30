package tableprinters

import (
	"fmt"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

const (
	dot = "●"
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
	case *apiv1.Tenant:
		return t.TenantTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Tenant:
		return t.TenantTable(d, wide)
	case *apiv1.IP:
		return t.IPTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.IP:
		return t.IPTable(d, wide)
	case *apiv1.Coupon:
		return t.CouponTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Coupon:
		return t.CouponTable(d, wide)
	case *apiv1.Cluster:
		return t.ClusterTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Cluster:
		return t.ClusterTable(d, wide)
	case *apiv1.ClusterStatusLastError:
		return t.ClusterStatusLastErrorTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.ClusterStatusLastError:
		return t.ClusterStatusLastErrorTable(d, wide)
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
	default:
		return nil, nil, fmt.Errorf("unknown table printer for type: %T", d)
	}
}
