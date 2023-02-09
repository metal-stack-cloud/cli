package tableprinters

import (
	"fmt"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
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
	case *apiv1.Tenant:
		return t.TenantTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.Tenant:
		return t.TenantTable(d, wide)
	case *apiv1.IP:
		return t.IPTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.IP:
		return t.IPTable(d, wide)
	default:
		return nil, nil, fmt.Errorf("unknown table printer for type: %T", d)
	}
}
