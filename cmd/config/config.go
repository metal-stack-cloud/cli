package config

import (
	"context"
	"io"

	adminv1client "github.com/metal-stack-cloud/api/go/client/admin/v1"
	apiv1client "github.com/metal-stack-cloud/api/go/client/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/printer"
)

type Config struct {
	Out           io.Writer
	Apiv1Client   apiv1client.Client
	Adminv1Client adminv1client.Client
	Ctx           context.Context
	Pf            *printer.PrinterFactory
}
