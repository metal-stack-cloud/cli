package e2erootcmd

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	client "github.com/metal-stack-cloud/api/go/client"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack-cloud/cli/cmd"
	"github.com/metal-stack-cloud/cli/cmd/completion"
	"github.com/metal-stack-cloud/cli/cmd/config"
	e2e_test "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type TestConfig struct {
	FsMocks     func(fs *afero.Afero)
	MockStdin   *bytes.Buffer
	ClientMocks *apitests.ClientMockFns
}

func NewRootCmd(t *testing.T, c *TestConfig) e2e_test.NewRootCmdFunc {
	return func() (*cobra.Command, *bytes.Buffer) {
		cl := client.New(&client.DialConfig{
			BaseURL:   "http://this-is-just-for-testing",
			UserAgent: "cli-test",
			Log:       slog.Default(),
		})

		if c.ClientMocks != nil {
			cl = apitests.New(t).Client(c.ClientMocks)
		}

		fs := afero.NewMemMapFs()
		if c.FsMocks != nil {
			c.FsMocks(&afero.Afero{
				Fs: fs,
			})
		}

		var in io.Reader
		if c.MockStdin != nil {
			in = bytes.NewReader(c.MockStdin.Bytes())
		}

		var out bytes.Buffer

		return cmd.NewRootCmd(&config.Config{
			Fs:        fs,
			Out:       &out,
			In:        in,
			PromptOut: io.Discard,
			Completion: &completion.Completion{
				Client: cl,
			},
			Client: cl,
		}), &out
	}
}
