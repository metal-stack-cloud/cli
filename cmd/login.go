package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
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
	loginCmd.Flags().String("context", "", "the context into which the token gets injected, if not specified it uses the current context or creates a context named default in case there is no current context set")

	genericcli.Must(loginCmd.RegisterFlagCompletionFunc("provider", c.Completion.LoginProviderCompletion))
	genericcli.Must(loginCmd.RegisterFlagCompletionFunc("context", c.ContextListCompletion))

	return loginCmd
}

func (l *login) login() error {
	provider := l.c.GetProvider()
	if provider == "" {
		return errors.New("provider must be specified")
	}

	// identify the context in which to inject the token

	ctxs, err := l.c.GetContexts()
	if err != nil {
		return err
	}

	ctxName := ctxs.CurrentContext
	if viper.IsSet("context") {
		ctxName = viper.GetString("context")
	}

	ctx, ok := ctxs.Get(ctxName)
	if !ok {
		newCtx := l.c.MustDefaultContext()
		newCtx.Name = "default"
		if viper.IsSet("context") {
			newCtx.Name = viper.GetString("context")
		}
		newCtx.ApiURL = pointer.Pointer(l.c.GetApiURL())

		ctxs.Contexts = append(ctxs.Contexts, &newCtx)

		ctx = &newCtx
	}

	ctx.Provider = provider

	// switch into new context
	ctxs.PreviousContext = ctxs.CurrentContext
	ctxs.CurrentContext = ctx.Name

	tokenChan := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tokenChan <- r.URL.Query().Get("token")

		http.Redirect(w, r, "https://metalstack.cloud", http.StatusSeeOther)
	})

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	server := http.Server{Addr: listener.Addr().String(), ReadTimeout: 2 * time.Second}

	go func() {
		if viper.GetBool("debug") {
			_, _ = fmt.Fprintf(l.c.Out, "Starting server at http://%s...\n", listener.Addr().String())
		}

		err = server.Serve(listener) //nolint
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("http server closed unexpectedly: %w", err))
		}
	}()

	url := fmt.Sprintf("%s/auth/%s?redirect-url=http://%s/callback", l.c.GetApiURL(), provider, listener.Addr().String())

	err = openBrowser(url)
	if err != nil {
		_, _ = fmt.Fprintf(l.c.Out, "the browser could not be opened, you can open it yourself on: %s", url)
		_, _ = fmt.Fprintf(l.c.Out, "the error was: %s", err.Error())
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
		mc := newApiClient(l.c.GetApiURL(), token)

		projects, err := mc.Apiv1().Project().List(context.Background(), connect.NewRequest(&apiv1.ProjectServiceListRequest{}))
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

	_, _ = fmt.Fprintf(l.c.Out, "%s login successful! Updated and activated context \"%s\"\n", color.GreenString("✔"), color.GreenString(ctx.Name))

	return nil
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	case "darwin":
		return exec.Command("open", url).Run()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
