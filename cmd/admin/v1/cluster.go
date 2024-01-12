package v1

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/kubernetes"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/afero"
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
		Description:     "manage cluster resources",
		Sorter:          sorters.ClusterSorter(),
		ValidArgsFn:     c.Completion.ClusterListCompletion,
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),

		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().BoolP("machines", "", false, "show machines of a cluster")
		},
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("purpose", "", "", "filter by purpose")
			// must(cmd.RegisterFlagCompletionFunc("id", c.Completion.ClusterListCompletion))

			genericcli.Must(cmd.RegisterFlagCompletionFunc("purpose", c.Completion.ClusterPurposeCompletion))
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

	kubeconfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")
	kubeconfigCmd.Flags().Bool("merge", true, "merges the kubeconfig into default kubeconfig instead of printing it to the console")
	kubeconfigCmd.Flags().String("kubeconfig", "", "specify an explicit path for the merged kubeconfig to be written, defaults to default kubeconfig paths if not provided")

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "fetch logs of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			req := &adminv1.ClusterServiceGetRequest{
				Uuid: id,
			}
			resp, err := c.Client.Adminv1().Cluster().Get(c.NewRequestContext(), connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to get cluster logs: %w", err)
			}
			return c.ListPrinter.Print(resp.Msg.Cluster.Status.LastErrors)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	monitoringCmd := &cobra.Command{
		Use:   "monitoring",
		Short: "fetch monitoring details of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}
			req := &adminv1.ClusterServiceGetRequest{
				Uuid: id,
			}

			resp, err := c.Client.Adminv1().Cluster().Get(c.NewRequestContext(), connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to get cluster monitoring: %w", err)
			}
			return c.DescribePrinter.Print(resp.Msg.Cluster.Monitoring)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	reconcileCmd := &cobra.Command{
		Use:   "reconcile",
		Short: "reconcile a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := genericcli.GetExactlyOneArg(args)
			if err != nil {
				return err
			}

			operation := adminv1.Operate_OPERATE_RECONCILE
			if viper.GetBool("maintain") {
				operation = adminv1.Operate_OPERATE_MAINTAIN
			}
			if viper.GetBool("retry") {
				operation = adminv1.Operate_OPERATE_RETRY
			}

			req := &adminv1.ClusterServiceOperateRequest{
				Uuid:    id,
				Operate: operation,
			}
			resp, err := c.Client.Adminv1().Cluster().Operate(c.NewRequestContext(), connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to reconcile cluster: %w", err)
			}

			return c.ListPrinter.Print(resp.Msg.Cluster)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	reconcileCmd.Flags().Bool("reconcile", true, "trigger cluster reconciliation")
	reconcileCmd.Flags().Bool("maintain", false, "trigger cluster maintain reconciliation")
	reconcileCmd.Flags().Bool("retry", false, "trigger cluster retry reconciliation")

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, reconcileCmd, logsCmd, monitoringCmd)
}

// TODO: implement firewall ssh, machine/firewall list

func (c *cluster) Create(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	showMachines := viper.GetBool("machines")
	req := &adminv1.ClusterServiceGetRequest{
		Uuid:         id,
		WithMachines: showMachines,
	}
	resp, err := c.c.Client.Adminv1().Cluster().Get(c.c.NewRequestContext(), connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}
	// FIXME refactor to use a Machine TablePrinter
	if showMachines {
		fmt.Println("Machines")
		fmt.Println()
		for _, m := range resp.Msg.Machines {
			fmt.Printf("%s: %s %s\n", m.Role, m.Uuid, m.Hostname)
		}
		fmt.Println()
	}

	c.c.ListPrinter.Print(resp.Msg.Cluster)
	return nil, nil
}

func (c *cluster) List() ([]*apiv1.Cluster, error) {
	// FIXME implement filters and paging

	req := &adminv1.ClusterServiceListRequest{}

	if viper.IsSet("purpose") {
		req.Purpose = pointer.Pointer(viper.GetString("purpose"))
	}

	resp, err := c.c.Client.Adminv1().Cluster().List(c.c.NewRequestContext(), connect.NewRequest(req))
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

	expiration := viper.GetDuration("expiration")
	req := &adminv1.ClusterServiceCredentialsRequest{
		Uuid:       id,
		Expiration: durationpb.New(expiration),
	}

	resp, err := c.c.Client.Adminv1().Cluster().Credentials(c.c.NewRequestContext(), connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to get cluster credentials: %w", err)
	}

	if !viper.GetBool("merge") {
		fmt.Fprintln(c.c.Out, resp.Msg.Kubeconfig)
		return nil
	}

	var (
		kubeconfigPath = viper.GetString("kubeconfig")
	)

	merged, err := kubernetes.MergeKubeconfig(c.c.Fs, []byte(resp.Msg.Kubeconfig), pointer.PointerOrNil(kubeconfigPath), nil) // FIXME: reverse lookup project name
	if err != nil {
		return err
	}

	err = afero.WriteFile(c.c.Fs, merged.Path, merged.Raw, 0600)
	if err != nil {
		return fmt.Errorf("unable to write merged kubeconfig: %w", err)
	}

	fmt.Fprintf(c.c.Out, "%s merged context %q into %s\n", color.GreenString("✔"), merged.ContextName, merged.Path)

	return nil
}
