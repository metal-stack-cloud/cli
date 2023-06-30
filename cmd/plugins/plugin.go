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

func AddPlugins(cmd *cobra.Command) error {
	// rootCmd.AddCommand(getCmd("/home/stefan/dev/github.com/metal-stack-cloud/cli-plugin/cli-plugin.so", "MainCmd"))
	h, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	pluginDir := path.Join(h, "."+config.ConfigDir, "plugins")
	if _, err := os.Stat(pluginDir); err != nil {
		if os.IsNotExist(err) {
			// no plugins
			return nil
		}
	}

	err = filepath.WalkDir(pluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if strings.HasSuffix(d.Name(), "-plugin.so") {
			cmdName, _, _ := strings.Cut(d.Name(), "-plugin.so")
			cmdName = cases.Title(language.English).String(cmdName)
			fmt.Printf("adding plugin from path:%q name:%q\n", path, cmdName)

			pluginCmd, err := getCmd(path, cmdName)
			if err != nil {
				return fmt.Errorf("unable to load plugin %q error %w", path, err)
			}
			cmd.AddCommand(pluginCmd)
		}
		return nil
	})

	return err
}

func getCmd(pluginPath, cmdName string) (*cobra.Command, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	b, err := p.Lookup(cmdName)
	if err != nil {
		return nil, err
	}
	f, err := p.Lookup("Init" + cmdName)
	if err == nil {
		f.(func())()
	}
	return *b.(**cobra.Command), nil
}
