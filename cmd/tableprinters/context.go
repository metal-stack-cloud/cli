package tableprinters

import (
	"github.com/fatih/color"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/viper"
)

func (t *TablePrinter) ContextTable(data *config.Contexts, wide bool) ([]string, [][]string, error) {
	var (
		header = []string{"", "Name", "Default Project"}
		rows   [][]string
	)

	if wide {
		header = []string{"", "Name", "Default Project", "API URL"}
	}

	for _, c := range data.Contexts {
		active := ""
		if c.Name == data.CurrentContext {
			active = color.GreenString("✔")
		}

		row := []string{active, c.Name, c.DefaultProject}
		if wide {
			url := pointer.SafeDeref(c.ApiURL)
			if url == "" {
				url = viper.GetString("api-url")
			}

			row = append(row, url)
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}
