package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/metal-stack-cloud/cli/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

// Execute is the entrypoint of the metal-go application
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
	name := "cli"
	rootCmd := &cobra.Command{
		Use:          name,
		Aliases:      []string{"m"},
		Short:        "a cli to manage metal-stack-cloud",
		Long:         "",
		SilenceUsage: true,
	}
	rootCmd.PersistentFlags().StringP("log-level", "l", "error", "configure loglevel, can be one of error|info|debug")
	rootCmd.PersistentFlags().StringP("config", "c", "", `alternative config file path, (default is ~/.cli/config.yaml).
Example config.yaml:

---
apitoken: "alongtoken"
...

`)
	lvl, err := rootCmd.PersistentFlags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	fmt.Printf("loglevel:%q", lvl)
	cfg := zap.NewProductionConfig()
	level, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		panic(err)
	}
	cfg.Level = level
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zlog, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	log := zlog.Sugar()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	apiClient, err := client.Dial(ctx, client.DialConfig{
		Endpoint:  "localhost:9090",
		Scheme:    client.GRPCS,
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtZXRhbC1zdGFjay1jbG91ZCIsInN1YiI6ImpvaG4uZG9lQGdpdGh1YiIsImV4cCI6MTY1MzIwMzQ0OSwibmJmIjoxNjUwNjExNDQ5LCJpYXQiOjE2NTA2MTE0NDksImp0aSI6Ijg4ZjFlMTBjLTk0YjItNDBkMS1iN2Q3LWZhYzQ4OGJiNmUxMiJ9.udyoOgQT1rZrrBbi8Io0tkj-W6Yuyu6Otz4NQQzI7qQ",
		Log:       log,
		UserAgent: "cli",
	})

	c := &config{
		client: apiClient,
		ctx:    ctx,
	}
	if err != nil {
		panic(err)
	}
	rootCmd.AddCommand(newVersionCmd(c))

	return rootCmd
}

type config struct {
	client client.Client
	ctx    context.Context
}
