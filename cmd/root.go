package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	client "github.com/metal-stack-cloud/api/go/client"
	adminv1client "github.com/metal-stack-cloud/api/go/client/admin/v1"
	apiv1client "github.com/metal-stack-cloud/api/go/client/api/v1"
	adminv1 "github.com/metal-stack-cloud/cli/cmd/admin/v1"
	apiv1 "github.com/metal-stack-cloud/cli/cmd/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

const (
	moduleName = "cli"
)

func Execute() {
	cmd := newRootCmd()

	err := cmd.Execute()
	if err != nil {
		if viper.GetBool("debug") {
			panic(err)
		}
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	config := &config.Config{
		Ctx: context.Background(),
	}

	rootCmd := &cobra.Command{
		Use:          moduleName,
		Aliases:      []string{"m"},
		Short:        "cli for managing entities in metal-stack-cloud",
		Long:         "",
		SilenceUsage: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			err := initConfig()
			if err != nil {
				return err
			}

			logger, err := newLogger()
			if err != nil {
				return err
			}

			apiclient, adminclient, err := newClient(logger)
			if err != nil {
				return err
			}

			config.Apiv1Client = apiclient
			config.Adminv1Client = adminclient

			config.Pf = &printer.PrinterFactory{Log: logger}
			config.Out = os.Stdout

			return nil
		},
	}
	rootCmd.PersistentFlags().StringP("log-level", "l", "error", "configure log level, can be one of error|info|debug")
	must(rootCmd.RegisterFlagCompletionFunc("log-level", cobra.FixedCompletions([]string{"error", "info", "debug"}, cobra.ShellCompDirectiveNoFileComp)))
	rootCmd.PersistentFlags().StringP("config", "c", "", `alternative config file path, (default is ~/.cli/config.yaml).
Example config.yaml:

---
apitoken: "alongtoken"
...
`)
	rootCmd.PersistentFlags().StringP("output-format", "o", "table", "output format (table|wide|markdown|json|yaml|template), wide is a table with more columns.")
	must(rootCmd.RegisterFlagCompletionFunc("output-format", cobra.FixedCompletions([]string{"table", "wide", "markdown", "json", "yaml", "template"}, cobra.ShellCompDirectiveNoFileComp)))
	rootCmd.PersistentFlags().StringP("template", "", "", `output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.`)
	rootCmd.PersistentFlags().Bool("force-color", false, "force colored output even without tty")

	rootCmd.PersistentFlags().String("api-url", "", "the url to the metal-stack cloud api")
	rootCmd.PersistentFlags().String("api-token", "", "the token used for api requests")
	rootCmd.PersistentFlags().String("api-ca-file", "", "the path to the ca file of the api server")

	must(viper.BindPFlags(rootCmd.PersistentFlags()))

	rootCmd.AddCommand(apiv1.NewVersionCmd(config))
	rootCmd.AddCommand(apiv1.NewHealthCmd(config))

	rootCmd.AddCommand(adminv1.NewCustomerCmd(config))

	return rootCmd
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func newLogger() (*zap.SugaredLogger, error) {
	lvl, err := zap.ParseAtomicLevel(viper.GetString("log-level"))
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zlog, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return zlog.Sugar(), nil
}

func initConfig() error {
	viper.SetEnvPrefix(strings.ToUpper(moduleName))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")
	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("config file path set explicitly, but unreadable: %w", err)
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s", moduleName))
		h, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("unable to figure out user home directory, skipping config lookup path: %v", err)
		} else {
			viper.AddConfigPath(fmt.Sprintf(h+"/.%s", moduleName))
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

func newClient(log *zap.SugaredLogger) (apiv1client.Client, adminv1client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	endpoint, err := url.Parse(viper.GetString("api-url"))
	if err != nil {
		return nil, nil, err
	}

	dialConfig := client.DialConfig{
		Endpoint: endpoint.Host,
		Token:    viper.GetString("api-token"),
		Credentials: &client.Credentials{
			ServerName: endpoint.Hostname(),
			CAFile:     viper.GetString("api-ca-file"),
		},
		Scheme:    client.GRPCS,
		UserAgent: "cli",
		Log:       log,
	}

	apiclient, err := apiv1client.New(ctx, dialConfig)
	if err != nil {
		return nil, nil, err
	}
	adminclient, err := adminv1client.New(ctx, dialConfig)
	if err != nil {
		return nil, nil, err
	}

	return apiclient, adminclient, nil
}
