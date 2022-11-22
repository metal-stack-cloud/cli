package v1

import (
	"fmt"

	v1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/v"
	"github.com/spf13/cobra"
)

type ip struct {
	Client string
	Server *v1.IP
}

func NewIpCmd(c *config.Config) *cobra.Command {
	ipCmd := &cobra.Command{
		Use: "ip",
		Short: "do something with ip",
		Long: "do something with ip",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: make it possible to allocate an IP
			i := ip{
				Client: v.V.String(),
			}
			// resp, err := c.Apiv1Client.IP().List(c.Ctx, &v1.IPServiceListRequest{})
			resp, err := c.Apiv1Client.IP().Get(c.Ctx, &v1.IPServiceGetRequest{})
			if err == nil {
				i.Server = resp.Ip
			}

			if err := c.Pf.NewPrinterDefaultYAML(c.Out).Print(i); err != nil {
				return err
			}

			if err != nil {
				return fmt.Errorf("failed to get server info: %w", err)
			}

			return nil
		},
	}
	return ipCmd
}