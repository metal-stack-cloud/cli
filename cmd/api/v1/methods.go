package v1

import (
	"fmt"
	"sort"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMethodsCmd(c *config.Config) *cobra.Command {
	methodCmd := &cobra.Command{
		Use:   "api-methods",
		Short: "show available api-methods of the metalstack.cloud api",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewRequestContext()
			defer cancel()

			if viper.GetBool("scoped") {

				req := &apiv1.MethodServiceTokenScopedListRequest{}

				resp, err := c.Client.Apiv1().Method().TokenScopedList(ctx, connect.NewRequest(req))
				if err != nil {
					return fmt.Errorf("failed to list methods: %w", err)
				}

				return printers.NewProtoYAMLPrinter().WithOut(c.Out).Print(resp.Msg)
			}

			var (
				methods []string
				req     = &apiv1.MethodServiceListRequest{}
			)

			resp, err := c.Client.Apiv1().Method().List(ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to list methods: %w", err)
			}

			methods = resp.Msg.GetMethods()

			sort.Strings(methods)

			for _, method := range methods {
				fmt.Fprintln(c.Out, method)
			}

			return nil
		},
	}

	methodCmd.Flags().Bool("scoped", false, "show accessible api-methods depending on the api access token")

	return methodCmd
}
