package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type login struct {
	c *config.Config
}

func newLoginCmd(c *config.Config) *cobra.Command {
	w := &login{
		c: c,
	}

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.login()
		},
	}

	loginCmd.Flags().String("provider", "", "the provider used to login with")
	genericcli.Must(loginCmd.RegisterFlagCompletionFunc("provider", cobra.FixedCompletions([]string{"github", "azure", "google"}, cobra.ShellCompDirectiveNoFileComp)))

	return loginCmd
}

func (l *login) login() error {
	provider := viper.GetString("provider")
	if provider == "" {
		return errors.New("provider must be specified")
	}

	ctxs, err := l.c.GetContexts()
	if err != nil {
		return err
	}

	ctx, ok := ctxs.Get(ctxs.CurrentContext)
	if !ok {
		defaultCtx := l.c.MustDefaultContext()
		defaultCtx.Name = "default"

		ctxs.PreviousContext = ctxs.CurrentContext
		ctxs.CurrentContext = ctx.Name
		ctxs.Contexts = append(ctxs.Contexts, &defaultCtx)

		ctx = &defaultCtx
	}

	tokenChan := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("error") != "" {
			http.Error(w, "Error: "+r.URL.Query().Get("error"), http.StatusBadRequest)
			return
		}

		tokenChan <- r.URL.Query().Get("token")
	})

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	server := http.Server{Addr: listener.Addr().String(), ReadTimeout: 2 * time.Second}

	go func() {
		fmt.Printf("Starting server at http://%s...\n", listener.Addr().String())
		err = server.Serve(listener) //nolint
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("http server closed unexpectedly: %w", err))
		}
	}()

	url := fmt.Sprintf("%s/auth/%s?redirect-url=http://%s/callback", l.c.GetApiURL(), provider, listener.Addr().String()) // TODO(vknabel): nicify please

	err = exec.Command("xdg-open", url).Run() //nolint
	if err != nil {
		return fmt.Errorf("error opening browser: %w", err)
	}

	token := <-tokenChan

	err = server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("unable to close http server: %w", err)
	}
	_ = listener.Close()

	if token == "" {
		return errors.New("no token was retrieved")
	}

	ctx.Token = token

	if ctx.DefaultProject == "" {
		projects, err := l.c.Client.Apiv1().Project().List(context.Background(), connect.NewRequest(&apiv1.ProjectServiceListRequest{}))
		if err != nil {
			return fmt.Errorf("unable to retrieve project list: %w", err)
		}

		idx := slices.IndexFunc(projects.Msg.Projects, func(p *apiv1.Project) bool {
			return p.IsDefaultProject
		})

		if idx >= 0 {
			ctx.DefaultProject = projects.Msg.Projects[idx].Uuid
		}
	}

	err = l.c.WriteContexts(ctxs)
	if err != nil {
		return err
	}

	fmt.Fprintf(l.c.Out, "%s login successful! Updated and activated context \"%s\"\n", color.GreenString("✔"), color.GreenString(ctx.Name))

	return nil
}
