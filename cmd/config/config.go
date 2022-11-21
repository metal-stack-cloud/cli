package config

import (
	"context"
	"io"

	adminv1client "github.com/metal-stack-cloud/api/go/client/admin/v1"
	apiv1client "github.com/metal-stack-cloud/api/go/client/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/printer"
	metalgo "github.com/metal-stack/metal-go"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metalctl/cmd/completion"
	"github.com/spf13/afero"
)

const (
	// BinaryName is the name of the cli in all help texts
	BinaryName = "cli"
	// ConfigDir is the directory in either the homedir or in /etc where the cli searches for a file config.yaml
	// also used as prefix for environment based configuration, e.g. METAL_STACK_CLOUD_ will be the variable prefix.
	ConfigDir = "metal-stack-cloud"
)

type Config struct {
	Fs            afero.Fs
	Out           io.Writer
	Comp		  *completion.Completion
	Apiv1Client   apiv1client.Client
	Adminv1Client adminv1client.Client
	Ctx           context.Context
	Pf            *printer.PrinterFactory
	Client         metalgo.Client
	DescribePrinter printers.Printer
	ListPrinter     printers.Printer
}


