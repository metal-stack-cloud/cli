package config

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/metal-stack-cloud/api/go/client"
	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
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
	ListPrinter     printers.Printer
	DescribePrinter printers.Printer
	Completion      *completion.Completion
	Context         Context
}

func (c *Config) NewRequestContext() (context.Context, context.CancelFunc) {
	timeout := c.Context.Timeout
	if timeout == nil {
		timeout = pointer.Pointer(30 * time.Second)
	}
	if viper.IsSet("timeout") {
		timeout = pointer.Pointer(viper.GetDuration("timeout"))
	}

	return context.WithTimeout(context.Background(), *timeout)
}

func HelpTemplate() string {
	return `Here is how an template configuration looks like:
~/.metal-stack-cloud/config.yaml
---
current: dev
previous: prod
contexts:
    - name: dev
    api-token: <dev-token>
    default-project: dev-project
    - name: prod
    api-token: <prod-token>
        default-project: prod-project
`
}

func DefaultConfigDirectory() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(h, "."+ConfigDir), nil
}

func ConfigPath() (string, error) {
	if viper.IsSet("config") {
		return viper.GetString("config"), nil
	}

	dir, err := DefaultConfigDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "config.yaml"), nil
}

func (c *Config) GetProject() string {
	if viper.IsSet("project") {
		return viper.GetString("project")
	}
	return c.Context.DefaultProject
}

func (c *Config) GetToken() string {
	if viper.IsSet("api-token") {
		return viper.GetString("api-token")
	}
	return c.Context.Token
}

func (c *Config) GetApiURL() string {
	if viper.IsSet("api-url") {
		return viper.GetString("api-url")
	}
	if c.Context.ApiURL != nil {
		return *c.Context.ApiURL
	}

	// fallback to the default specified by viper
	return viper.GetString("api-url")
}
