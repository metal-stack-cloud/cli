package cmd

import (
	"fmt"

	v1 "github.com/metal-stack-cloud/api-server/api/v1"
	"github.com/metal-stack-cloud/cli/api"
	"github.com/metal-stack/v"
	"github.com/spf13/cobra"
)

func newVersionCmd(c *config) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print the client and server version information",
		Long:  "print the client and server version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			v := api.Version{
				Client: v.V.String(),
			}
			resp, err := c.client.Version().Get(c.ctx, &v1.VersionServiceGetRequest{})
			if err == nil {
				v.Server = resp.Version.String()
			}
			fmt.Print(v)
			if err != nil {
				return fmt.Errorf("failed to get server info: %w", err)
			}
			return nil
		},
	}
	return versionCmd
}
