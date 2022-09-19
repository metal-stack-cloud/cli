package printer

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type PrinterFactory struct {
	Log *zap.SugaredLogger
}

func (p *PrinterFactory) NewPrinter(out io.Writer) printers.Printer {
	var printer printers.Printer

	switch format := viper.GetString("output-format"); format {
	case "yaml":
		printer = printers.NewProtoYAMLPrinter().WithFallback(true)
	case "json":
		printer = printers.NewProtoJSONPrinter().WithFallback(true)
	case "table", "wide", "markdown":
		tp := New()
		cfg := &printers.TablePrinterConfig{
			ToHeaderAndRows: tp.ToHeaderAndRows,
			Wide:            format == "wide",
			Markdown:        format == "markdown",
			NoHeaders:       viper.GetBool("no-headers"),
		}
		tablePrinter := printers.NewTablePrinter(cfg).WithOut(out) // FIXME configurable
		tp.SetPrinter(tablePrinter)
		tp.SetLastEventErrorThreshold(viper.GetDuration("last-event-error-threshold"))
		printer = tablePrinter
	case "template":
		printer = printers.NewTemplatePrinter(viper.GetString("template"))
	default:
		p.Log.Fatalf("unknown output format: %q", format)
	}

	if viper.IsSet("force-color") {
		enabled := viper.GetBool("force-color")
		if enabled {
			color.NoColor = false
		} else {
			color.NoColor = true
		}
	}

	return printer
}
func (p *PrinterFactory) NewPrinterDefaultYAML(out io.Writer) printers.Printer {
	if viper.IsSet("output-format") {
		return p.NewPrinter(out)
	}
	return printers.NewProtoYAMLPrinter().WithOut(out).WithFallback(true)
}

type TablePrinter struct {
	t                       *printers.TablePrinter
	lastEventErrorThreshold time.Duration
}

func New() *TablePrinter {
	return &TablePrinter{}
}

func (t *TablePrinter) SetPrinter(printer *printers.TablePrinter) {
	t.t = printer
}

func (t *TablePrinter) SetLastEventErrorThreshold(threshold time.Duration) {
	t.lastEventErrorThreshold = threshold
}

func (t *TablePrinter) ToHeaderAndRows(data any, wide bool) ([]string, [][]string, error) {
	switch d := data.(type) {
	case *apiv1.PaymentCustomer:
		return t.CustomerTable(pointer.WrapInSlice(d), wide)
	case []*apiv1.PaymentCustomer:
		return t.CustomerTable(d, wide)
	default:
		return nil, nil, fmt.Errorf("unknown table printer for type: %T", d)
	}
}
