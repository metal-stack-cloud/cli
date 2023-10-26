package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	client "github.com/metal-stack-cloud/api/go/client"

	adminv1 "github.com/metal-stack-cloud/cli/cmd/admin/v1"
	apiv1 "github.com/metal-stack-cloud/cli/cmd/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

func Execute() {
	cfg := &config.Config{
		Ctx:        context.Background(),
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
		fmt.Printf("%+v\n", err)
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
			must(viper.BindPFlags(cmd.Flags()))
			must(viper.BindPFlags(cmd.PersistentFlags()))

			must(readConfigFile())
			initConfigWithViperCtx(c)

			return nil
		},
	}
	rootCmd.PersistentFlags().StringP("config", "c", "", `alternative config file path, (default is ~/.metal-stack-cloud/config.yaml).
Example config.yaml:

---
apitoken: "alongtoken"
...
`)
	rootCmd.PersistentFlags().StringP("output-format", "o", "table", "output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact.")
	must(rootCmd.RegisterFlagCompletionFunc("output-format", cobra.FixedCompletions([]string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp)))
	rootCmd.PersistentFlags().StringP("template", "", "", `output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.`)
	rootCmd.PersistentFlags().Bool("force-color", false, "force colored output even without tty")
	rootCmd.PersistentFlags().Bool("debug", false, "debug output")

	rootCmd.PersistentFlags().String("api-url", "", "the url to the metalstack.cloud api")
	rootCmd.PersistentFlags().String("api-token", "", "the token used for api requests")
	rootCmd.PersistentFlags().String("api-ca-file", "", "the path to the ca file of the api server")

	must(viper.BindPFlags(rootCmd.PersistentFlags()))

	adminv1.AddCmds(rootCmd, c)
	apiv1.AddCmds(rootCmd, c)

	return rootCmd
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func readConfigFile() error {
	viper.SetEnvPrefix(strings.ToUpper(config.ConfigDir))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")
	cfgFile := viper.GetString("config")

	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("config file path set explicitly, but unreadable: %w", err)
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s", config.ConfigDir))
		h, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("unable to figure out user home directory, skipping config lookup path: %v", err)
		} else {
			viper.AddConfigPath(fmt.Sprintf(h+"/.%s", config.ConfigDir))
		}
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			usedCfg := viper.ConfigFileUsed()
			if usedCfg != "" {
				return fmt.Errorf("config %s file unreadable: %w", usedCfg, err)
			}
		}
	}

	return nil
}

func initConfigWithViperCtx(c *config.Config) {
	c.Log = newLogger()

	c.ListPrinter = newPrinterFromCLI(c.Log, c.Out)
	c.DescribePrinter = defaultToYAMLPrinter(c.Log, c.Out)

	if c.Client != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dialConfig := client.DialConfig{
		BaseURL:   viper.GetString("api-url"),
		Token:     viper.GetString("api-token"),
		UserAgent: "metal-stack-cloud-cli",
		Debug:     viper.GetBool("debug"),
	}

	mc := client.New(dialConfig)

	c.Client = mc
	c.Completion.Client = mc
	c.Completion.Ctx = ctx
}

func newLogger() *zap.SugaredLogger {
	level := zapcore.ErrorLevel
	if viper.GetBool("debug") {
		level = zapcore.DebugLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zlog, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return zlog.Sugar()
}
