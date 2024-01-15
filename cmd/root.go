package cmd

import (
	"os"

	client "github.com/metal-stack-cloud/api/go/client"
	"github.com/metal-stack/metal-lib/pkg/genericcli"

	adminv1 "github.com/metal-stack-cloud/cli/cmd/admin/v1"
	apiv1 "github.com/metal-stack-cloud/cli/cmd/api/v1"

	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func Execute() {
	cfg := &config.Config{
		Fs:         afero.NewOsFs(),
		Out:        os.Stdout,
		Completion: &completion.Completion{},
	}

	cmd := newRootCmd(cfg)

	err := cmd.Execute()
	if err != nil {
		if viper.GetBool("debug") {
			panic(err)
		}

		os.Exit(1)
	}
}

func newRootCmd(c *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          config.BinaryName,
		Aliases:      []string{"m"},
		Short:        "cli for managing entities in metal-stack-cloud",
		Long:         "",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viper.SetFs(c.Fs)

			genericcli.Must(viper.BindPFlags(cmd.Flags()))
			genericcli.Must(viper.BindPFlags(cmd.PersistentFlags()))

			return initConfigWithViperCtx(c)
		},
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", "alternative config file path, (default is ~/.metal-stack-cloud/config.yaml)")
	rootCmd.PersistentFlags().StringP("output-format", "o", "table", "output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact.")
	genericcli.Must(rootCmd.RegisterFlagCompletionFunc("output-format", cobra.FixedCompletions([]string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp)))
	rootCmd.PersistentFlags().StringP("template", "", "", `output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.`)
	rootCmd.PersistentFlags().Bool("force-color", false, "force colored output even without tty")
	rootCmd.PersistentFlags().Bool("debug", false, "debug output")
	rootCmd.PersistentFlags().Duration("timeout", 0, "request timeout used for api requests")

	rootCmd.PersistentFlags().String("api-url", "https://api.metalstack.cloud", "the url to the metalstack.cloud api")
	rootCmd.PersistentFlags().String("api-token", "", "the token used for api requests")
	rootCmd.PersistentFlags().String("api-ca-file", "", "the path to the ca file of the api server")

	genericcli.Must(viper.BindPFlags(rootCmd.PersistentFlags()))

	rootCmd.AddCommand(newContextCmd(c))
	adminv1.AddCmds(rootCmd, c)
	apiv1.AddCmds(rootCmd, c)

	return rootCmd
}

func initConfigWithViperCtx(c *config.Config) error {
	c.Context = c.MustDefaultContext()

	listPrinter, err := newPrinterFromCLI(c.Out)
	if err != nil {
		return err
	}
	describePrinter, err := defaultToYAMLPrinter(c.Out)
	if err != nil {
		return err
	}

	c.ListPrinter = listPrinter
	c.DescribePrinter = describePrinter

	if c.Client != nil {
		return nil
	}

	dialConfig := client.DialConfig{
		BaseURL:   c.GetApiURL(),
		Token:     c.GetToken(),
		UserAgent: "metal-stack-cloud-cli",
		Debug:     viper.GetBool("debug"),
	}

	mc := client.New(dialConfig)

	c.Client = mc
	c.Completion.Client = mc
	c.Completion.Ctx = context.Background()
	c.Completion.Project = c.GetProject()

	return nil
}
