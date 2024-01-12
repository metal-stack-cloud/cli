package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// Contexts contains all configuration contexts
type Contexts struct {
	CurrentContext  string     `json:"current-context"`
	PreviousContext string     `json:"previous-context"`
	Contexts        []*Context `json:"contexts"`
}

// Context configure
type Context struct {
	Name           string         `json:"name"`
	ApiURL         *string        `json:"api-url,omitempty"`
	Token          string         `json:"api-token"`
	DefaultProject string         `json:"default-project"`
	Timeout        *time.Duration `json:"timeout,omitempty"`
}

func defaultCtx() Context {
	return Context{
		ApiURL: pointer.PointerOrNil(viper.GetString("api-url")),
		Token:  viper.GetString("api-token"),
	}
}

func (cs Contexts) GetContext(name string) (*Context, bool) {
	for _, context := range cs.Contexts {
		context := context

		if context.Name == name {
			return context, true
		}
	}

	return nil, false
}

func (cs Contexts) Validate() error {
	names := map[string]bool{}
	for _, context := range cs.Contexts {
		context := context

		names[context.Name] = true
	}

	if len(cs.Contexts) != len(names) {
		return fmt.Errorf("context names must be unique")
	}

	return nil
}

func (cs *Contexts) Delete(name string) {
	cs.Contexts = slices.DeleteFunc(cs.Contexts, func(ctx *Context) bool {
		return ctx.Name == name
	})
}

func GetContexts() (*Contexts, error) {
	var ctxs Contexts
	path := viper.ConfigFileUsed()

	c, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Contexts{}, nil
		}

		return nil, fmt.Errorf("unable to read config.yaml: %w", err)
	}

	err = yaml.Unmarshal(c, &ctxs)
	return &ctxs, err
}

func WriteContexts(ctxs *Contexts) error {
	if err := ctxs.Validate(); err != nil {
		return err
	}

	c, err := yaml.Marshal(ctxs)
	if err != nil {
		return err
	}

	path := viper.ConfigFileUsed()
	if path == "" {
		path, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, c, 0600)
	if err != nil {
		return err
	}

	return nil
}

func MustDefaultContext() Context {
	ctxs, err := GetContexts()
	if err != nil {
		return defaultCtx()
	}
	ctx, ok := ctxs.GetContext(ctxs.CurrentContext)
	if !ok {
		return defaultCtx()
	}
	return *ctx
}

func ContextListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ctxs, err := GetContexts()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, ctx := range ctxs.Contexts {
		ctx := ctx
		names = append(names, ctx.Name)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c Context) GetProject() string {
	if viper.IsSet("project") {
		return viper.GetString("project")
	}
	return c.DefaultProject
}

func (c Context) GetToken() string {
	if viper.IsSet("api-token") {
		return viper.GetString("api-token")
	}
	return c.Token
}

func (c Context) GetApiURL() string {
	if viper.IsSet("api-url") {
		return viper.GetString("api-url")
	}
	if c.ApiURL != nil {
		return *c.ApiURL
	}
	return viper.GetString("api-url")
}
