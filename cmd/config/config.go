package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"connectrpc.com/connect"
	"github.com/metal-stack-cloud/api/go/client"
	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
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
	In              io.Reader
	Out             io.Writer
	PromptOut       io.Writer
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

func (c *Config) GetTenant() (string, error) {
	if viper.IsSet("tenant") {
		return viper.GetString("tenant"), nil
	}

	if c.GetProject() != "" {
		ctx, cancel := c.NewRequestContext()
		defer cancel()

		projectResp, err := c.Client.Apiv1().Project().Get(ctx, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
			Project: c.GetProject(),
		}))
		if err != nil {
			return "", fmt.Errorf("unable to derive tenant from project: %w", err)
		}

		return projectResp.Msg.Project.Tenant, nil
	}

	ctx, cancel := c.NewRequestContext()
	defer cancel()

	user, err := c.Client.Apiv1().User().Get(ctx, connect.NewRequest(&apiv1.UserServiceGetRequest{}))
	if err != nil {
		return "", fmt.Errorf("unable to derive tenant from token: %w", err)
	}

	return user.Msg.User.Login, nil
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

func (c *Config) GetProvider() string {
	if viper.IsSet("provider") {
		return viper.GetString("provider")
	}
	return c.Context.Provider
}
