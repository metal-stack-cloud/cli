package cmd

import (
	"errors"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"

	client "github.com/metal-stack-cloud/api/go/client"
	"github.com/metal-stack/metal-lib/pkg/genericcli"

	adminv1cmds "github.com/metal-stack-cloud/cli/cmd/admin/v1"
	apiv1cmds "github.com/metal-stack-cloud/cli/cmd/api/v1"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"

	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func Execute() {
	cfg := &config.Config{
		Fs:         afero.NewOsFs(),
		Out:        os.Stdout,
		PromptOut:  os.Stdout,
		In:         os.Stdin,
		Completion: &completion.Completion{},
	}

	cmd := newRootCmd(cfg)
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err != nil {
		printErrorWithHints(cfg, cmd, err)

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

	genericcli.Must(viper.BindPFlags(rootCmd.PersistentFlags()))

	markdownCmd := &cobra.Command{
		Use:   "markdown",
		Short: "create markdown documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, "./docs")
		},
		DisableAutoGenTag: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			recursiveAutoGenDisable(rootCmd)
		},
	}

	rootCmd.AddCommand(newContextCmd(c), markdownCmd, newLoginCmd(c))
	adminv1cmds.AddCmds(rootCmd, c)
	apiv1cmds.AddCmds(rootCmd, c)

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

	mc := newApiClient(c.GetApiURL(), c.GetToken())

	c.Client = mc
	c.Completion.Client = mc
	c.Completion.Ctx = context.Background()
	c.Completion.Project = c.GetProject()

	return nil
}

func newApiClient(apiURL, token string) client.Client {
	dialConfig := client.DialConfig{
		BaseURL:   apiURL,
		Token:     token,
		UserAgent: "metal-stack-cloud-cli",
		Debug:     viper.GetBool("debug"),
	}

	return client.New(dialConfig)
}

func recursiveAutoGenDisable(cmd *cobra.Command) {
	cmd.DisableAutoGenTag = true
	for _, child := range cmd.Commands() {
		recursiveAutoGenDisable(child)
	}
}

func printErrorWithHints(c *config.Config, cmd *cobra.Command, err error) {
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) || connectErr.Code() != connect.CodeUnauthenticated {
		cmd.PrintErrln(err)
		return
	}

	token := c.GetToken()
	if token != "" {
		type claims struct {
			jwt.RegisteredClaims
			Type string `json:"type"`
		}

		cs := &claims{}
		_, _, err := new(jwt.Parser).ParseUnverified(string(token), cs)
		if err == nil && cs.ExpiresAt != nil {
			if cs.ExpiresAt.Before(time.Now()) {
				switch cs.Type {
				case apiv1.TokenType_TOKEN_TYPE_API.String():

					cmd.PrintErrf("The API token has expired at %s.\nCreate a new API token or use the login command.\n", cs.ExpiresAt.String())
				case apiv1.TokenType_TOKEN_TYPE_CONSOLE.String(), apiv1.TokenType_TOKEN_TYPE_UNSPECIFIED.String():
					fallthrough
				default:
					cmd.PrintErrf("The token has expired at %s.\nPlease use the login command.\n", cs.ExpiresAt.String())
				}
			}
		}
	}
}
