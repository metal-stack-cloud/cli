package config

import (
	"fmt"
	"os"

	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// Contexts contains all configuration contexts of metalctl
type Contexts struct {
	CurrentContext  string     `json:"current-context"`
	PreviousContext string     `json:"previous-context"`
	Contexts        []*Context `json:"contexts"`
}

// Context configure metalctl behaviour
type Context struct {
	Name           string  `json:"name"`
	ApiURL         *string `json:"api-url,omitempty"`
	Token          string  `json:"api-token"`
	DefaultProject string  `json:"default-project"`
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

func GetContexts() (*Contexts, error) {
	var ctxs Contexts
	cfgFile := viper.ConfigFileUsed()
	c, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read config, please create a config.yaml")
	}
	err = yaml.Unmarshal(c, &ctxs)
	return &ctxs, err
}

func WriteContexts(ctxs *Contexts) error {
	c, err := yaml.Marshal(ctxs)
	if err != nil {
		return err
	}
	cfgFile := viper.ConfigFileUsed()
	err = os.WriteFile(cfgFile, c, 0600)
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
	return ""
}
