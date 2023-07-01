package config

import (
	"context"
	"io"

	adminv1client "github.com/metal-stack-cloud/api/go/client/admin/v1"
	apiv1client "github.com/metal-stack-cloud/api/go/client/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

const (
	// BinaryName is the name of the cli in all help texts
	BinaryName = "metal"
	// ConfigDir is the directory in either the homedir or in /etc where the cli searches for a file config.yaml
	// also used as prefix for environment based configuration, e.g. METAL_STACK_CLOUD_ will be the variable prefix.
	ConfigDir    = "metal-stack-cloud"
	PluginDir    = "plugins"
	PluginSuffix = "-plugin.so"
)

type Config struct {
	Fs              afero.Fs
	Out             io.Writer
	Apiv1Client     apiv1client.Client
	Adminv1Client   adminv1client.Client
	Ctx             context.Context
	Log             *zap.SugaredLogger
	ListPrinter     printers.Printer
	DescribePrinter printers.Printer
}
