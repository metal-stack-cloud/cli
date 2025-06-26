package cmd

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/metal-stack-cloud/cli/cmd/tableprinters"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/viper"
)

func newPrinterFromCLI(out io.Writer) (printers.Printer, error) {
	var printer printers.Printer

	switch format := viper.GetString("output-format"); format {
	case "yaml":
		printer = printers.NewProtoYAMLPrinter().WithFallback(true).WithOut(out)
	case "json":
		printer = printers.NewProtoJSONPrinter().WithFallback(true).WithOut(out)
	case "yamlraw":
		printer = printers.NewYAMLPrinter().WithOut(out)
	case "jsonraw":
		printer = printers.NewJSONPrinter().WithOut(out)
	case "table", "wide", "markdown":
		tp := tableprinters.New()
		cfg := &printers.TablePrinterConfig{
			ToHeaderAndRows: tp.ToHeaderAndRows,
			Wide:            format == "wide",
			Markdown:        format == "markdown",
			NoHeaders:       viper.GetBool("no-headers"),
		}
		tablePrinter := printers.NewTablePrinter(cfg).WithOut(out)
		tp.SetPrinter(tablePrinter)
		printer = tablePrinter
	case "template":
		printer = printers.NewTemplatePrinter(viper.GetString("template")).WithOut(out)
	default:
		return nil, fmt.Errorf("unknown output format: %q", format)
	}

	if viper.IsSet("force-color") {
		enabled := viper.GetBool("force-color")
		if enabled {
			color.NoColor = false
		} else {
			color.NoColor = true
		}
	}

	return printer, nil
}

func defaultToYAMLPrinter(out io.Writer) (printers.Printer, error) {
	if viper.IsSet("output-format") {
		return newPrinterFromCLI(out)
	}
	return printers.NewProtoYAMLPrinter().WithFallback(true).WithOut(out), nil
}
