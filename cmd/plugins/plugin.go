package plugins

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/cobra"
)

func NewPluginCommand() *cobra.Command {
	pluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "print the configured plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := getPlugins(nil)
			if err != nil {
				return err
			}
			if len(ps) < 1 {
				fmt.Println("no plugins available")
			}
			for _, p := range ps {
				fmt.Printf("%s:%q version: %q\n", p.name, p.path, p.version)
			}
			return nil
		},
	}
	return pluginCmd
}

func AddPlugins(cmd *cobra.Command, cfg *config.Config) error {
	ps, err := getPlugins(cfg)
	if err != nil {
		return err
	}
	for _, p := range ps {
		cmd.AddCommand(p.command)
	}
	return nil
}

type cliplugin struct {
	name    string
	command *cobra.Command
	path    string
	version string
}

func getPlugins(cfg *config.Config) ([]cliplugin, error) {
	var ps []cliplugin

	h, err := os.UserHomeDir()
	if err != nil {
		return ps, err
	}

	pluginDir := path.Join(h, "."+config.ConfigDir, config.PluginDir)
	if _, err := os.Stat(pluginDir); err != nil {
		if os.IsNotExist(err) {
			// no plugins
			return ps, nil
		}
	}

	err = filepath.WalkDir(pluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(d.Name(), config.PluginSuffix) {
			return nil
		}

		cmdName, _, _ := strings.Cut(d.Name(), config.PluginSuffix)
		cmdName = cases.Title(language.English).String(cmdName)

		pluginCmd, version, err := getCmd(path, cmdName, cfg)
		if err != nil {
			return fmt.Errorf("unable to load plugin %q error %w", path, err)
		}
		ps = append(ps, cliplugin{
			name:    cmdName,
			command: pluginCmd,
			path:    path,
			version: version,
		})
		return nil
	})

	return ps, err
}

func getCmd(pluginPath, cmdName string, cfg *config.Config) (*cobra.Command, string, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, "unknown", err
	}
	b, err := p.Lookup(cmdName)
	if err != nil {
		return nil, "unknown", err
	}
	c, err := p.Lookup("Config")
	if err != nil {
		return nil, "unknown", err
	}
	if cfg != nil {
		*c.(*config.Config) = *cfg
	}

	v, err := p.Lookup("V")
	if err != nil {
		return nil, "unknown", err
	}

	version := *v.(*string)

	return *b.(**cobra.Command), version, nil
}
