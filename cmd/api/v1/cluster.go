package v1

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type cluster struct {
	c *config.Config
}

func newClusterCmd(c *config.Config) *cobra.Command {
	w := &cluster{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, any, *apiv1.Cluster]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[any, any, *apiv1.Cluster](w).WithFS(c.Fs),
		Singular:        "cluster",
		Plural:          "clusters",
		Description:     "manage kubernetes clusters",
		Sorter:          sorters.ClusterSorter(),
		ValidArgsFn:     c.Completion.ClusterListCompletion,
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project for which to list clusters")
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("project", "", "the project for which to describe the cluster")
		},
	}

	kubeconfigCmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "fetch kubeconfig of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}

			if !viper.IsSet("project") {
				return fmt.Errorf("project is required to be set")
			}

			req := &apiv1.ClusterServiceGetCredentialsRequest{
				Uuid:       id,
				Project:    viper.GetString("project"),
				Expiration: durationpb.New(viper.GetDuration("expiration")),
			}

			resp, err := c.Client.Apiv1().Cluster().GetCredentials(c.Ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to get cluster credentials: %w", err)
			}

			fmt.Fprintln(w.c.Out, resp.Msg.Kubeconfig)

			return nil
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	kubeconfigCmd.Flags().String("project", "", "the project in which the cluster resides for which to get the kubeconfig for")
	kubeconfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")

	monitoringCmd := &cobra.Command{
		Use:   "monitoring",
		Short: "fetch endpoints and access credentials to cluster monitoring",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}

			cluster, err := w.Get(id)
			if err != nil {
				return err
			}

			return c.DescribePrinter.Print(cluster.Monitoring)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, monitoringCmd)
}

func (c *cluster) Create(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	if !viper.IsSet("project") {
		return nil, fmt.Errorf("project is required to be set")
	}

	req := &apiv1.ClusterServiceGetRequest{
		Uuid:    id,
		Project: viper.GetString("project"),
	}

	resp, err := c.c.Client.Apiv1().Cluster().Get(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) List() ([]*apiv1.Cluster, error) {
	if !viper.IsSet("project") {
		return nil, fmt.Errorf("project is required to be set")
	}

	req := &apiv1.ClusterServiceListRequest{
		Project: viper.GetString("project"),
	}

	resp, err := c.c.Client.Apiv1().Cluster().List(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return resp.Msg.Clusters, nil
}

func (c *cluster) Convert(r *apiv1.Cluster) (string, any, any, error) {
	panic("unimplemented")
}

func (c *cluster) Update(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}
