package config

import (
	"context"
	"io"

	adminv1client "github.com/metal-stack-cloud/api/go/client/admin/v1"
	apiv1client "github.com/metal-stack-cloud/api/go/client/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/printer"
	"github.com/spf13/afero"
)

const (
	BinaryName = "cli"
)

type Config struct {
	Fs            afero.Fs
	Out           io.Writer
	Apiv1Client   apiv1client.Client
	Adminv1Client adminv1client.Client
	Ctx           context.Context
	Pf            *printer.PrinterFactory
}
