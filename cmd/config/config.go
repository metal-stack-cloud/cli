package config

import (
	"context"
	"io"

	"github.com/metal-stack-cloud/api/go/client"
	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/afero"
)

const (
	// BinaryName is the name of the cli in all help texts
	BinaryName = "metal"
	// ConfigDir is the directory in either the homedir or in /etc where the cli searches for a file config.yaml
	// also used as prefix for environment based configuration, e.g. METAL_STACK_CLOUD_ will be the variable prefix.
	ConfigDir = "metal-stack-cloud"
)

type Config struct {
	Fs              afero.Fs
	Out             io.Writer
	Client          client.Client
	Ctx             context.Context
	ListPrinter     printers.Printer
	DescribePrinter printers.Printer
	Completion      *completion.Completion
}
