package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	client "github.com/metal-stack-cloud/api/go/client"
	apiv1client "github.com/metal-stack-cloud/api/go/client/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

const (
	moduleName = "cli"
)

type config struct {
	apiv1client apiv1client.Client
	ctx         context.Context
	pf          *printerFactory
}

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
	config := &config{
		ctx: context.Background(),
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

			client, err := newClient(logger)
			if err != nil {
				return err
			}

			config.apiv1client = client
			config.pf = &printerFactory{log: logger}

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

	rootCmd.AddCommand(newVersionCmd(config))
	rootCmd.AddCommand(newHealthCmd(config))

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

func newClient(log *zap.SugaredLogger) (apiv1client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	endpoint, err := url.Parse(viper.GetString("api-url"))
	if err != nil {
		return nil, err
	}

	c, err := apiv1client.New(ctx, client.DialConfig{
		Endpoint: endpoint.Host,
		Token:    viper.GetString("api-token"),
		Credentials: &client.Credentials{
			ServerName: endpoint.Hostname(),
			CAFile:     viper.GetString("api-ca-file"),
		},
		Scheme:    client.GRPCS,
		UserAgent: "cli",
		Log:       log,
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

type printerFactory struct {
	log *zap.SugaredLogger
}

func (p *printerFactory) newPrinter() printers.Printer {
	var printer printers.Printer

	switch format := viper.GetString("output-format"); format {
	case "yaml":
		printer = printers.NewProtoYAMLPrinter()
	case "json":
		printer = printers.NewProtoJSONPrinter()
	case "template":
		printer = printers.NewTemplatePrinter(viper.GetString("template"))
	default:
		p.log.Fatalf("unknown output format: %q", format)
	}

	if viper.IsSet("force-color") {
		enabled := viper.GetBool("force-color")
		if enabled {
			color.NoColor = false
		} else {
			color.NoColor = true
		}
	}

	return printer
}
func (p *printerFactory) newPrinterDefaultYAML() printers.Printer {
	if viper.IsSet("output-format") {
		return p.newPrinter()
	}
	return printers.NewProtoYAMLPrinter()
}
