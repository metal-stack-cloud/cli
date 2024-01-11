package v1

import (
	"fmt"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/kubernetes"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
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
			return w.kubeconfig(args)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	kubeconfigCmd.Flags().String("project", "", "the project in which the cluster resides for which to get the kubeconfig for")
	kubeconfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")
	kubeconfigCmd.Flags().Bool("merge", true, "merges the kubeconfig into default kubeconfig instead of printing it to the console")
	kubeconfigCmd.Flags().String("kubeconfig", "", "specify an explicit path for the merged kubeconfig to be written, defaults to default kubeconfig paths if not provided")

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
	req := &apiv1.ClusterServiceGetRequest{
		Uuid:    id,
		Project: c.c.Context.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Cluster().Get(c.c.NewRequestContext(), connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) List() ([]*apiv1.Cluster, error) {
	req := &apiv1.ClusterServiceListRequest{
		Project: c.c.Context.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Cluster().List(c.c.NewRequestContext(), connect.NewRequest(req))
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

func (c *cluster) kubeconfig(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	req := &apiv1.ClusterServiceGetCredentialsRequest{
		Uuid:       id,
		Project:    c.c.Context.GetProject(),
		Expiration: durationpb.New(viper.GetDuration("expiration")),
	}

	resp, err := c.c.Client.Apiv1().Cluster().GetCredentials(c.c.NewRequestContext(), connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to get cluster credentials: %w", err)
	}

	if !viper.GetBool("merge") {
		fmt.Fprintln(c.c.Out, resp.Msg.Kubeconfig)
		return nil
	}

	var (
		kubeconfigPath = viper.GetString("kubeconfig")
		projectName    = c.c.Context.GetProject() // FIXME: reverse lookup project name from
	)

	merged, err := kubernetes.MergeKubeconfig([]byte(resp.Msg.Kubeconfig), pointer.PointerOrNil(kubeconfigPath), &projectName)
	if err != nil {
		return err
	}

	err = os.WriteFile(merged.Path, merged.Raw, 0600)
	if err != nil {
		return fmt.Errorf("unable to write merged kubeconfig: %w", err)
	}

	fmt.Fprintf(c.c.Out, "%s merged context %q into %s\n", color.GreenString("✔"), merged.ContextName, merged.Path)

	return nil
}
