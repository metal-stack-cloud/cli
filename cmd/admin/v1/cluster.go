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
	"github.com/metal-stack/metal-lib/pkg/tag"
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
		ValidArgsFn:     c.Completion.ClusterAdminListCompletion,
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("id", "", "filter by id")
			cmd.Flags().StringP("project", "p", "", "filter by project")
			cmd.Flags().String("tenant", "", "filter by tenant")
			cmd.Flags().String("partition", "", "filter by partition")
			cmd.Flags().String("seed", "", "filter by seed")
			cmd.Flags().String("name", "", "filter by name")
			cmd.Flags().String("labels", "", "filter by labels")
			cmd.Flags().String("purpose", "", "filter by purpose")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.Completion.ClusterAdminListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("name", c.Completion.ClusterNameAdminListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("purpose", c.Completion.ClusterPurposeCompletion))
		},
	}

	// metal admin cluster kubeconfig

	kubeconfigCmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "fetch kubeconfig of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.kubeconfig(args)
		},
		ValidArgsFunction: c.Completion.ClusterAdminListCompletion,
	}

	kubeconfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")
	kubeconfigCmd.Flags().Bool("merge", true, "merges the kubeconfig into default kubeconfig instead of printing it to the console")
	kubeconfigCmd.Flags().String("kubeconfig", "", "specify an explicit path for the merged kubeconfig to be written, defaults to default kubeconfig paths if not provided")

	// metal admin cluster logs

	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "fetch logs of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.logs(args)
		},
		ValidArgsFunction: c.Completion.ClusterAdminListCompletion,
	}

	// metal admin cluster reconcile

	reconcileCmd := &cobra.Command{
		Use:   "reconcile",
		Short: "reconcile a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.reconcile(args)
		},
		ValidArgsFunction: c.Completion.ClusterAdminListCompletion,
	}

	reconcileCmd.Flags().String("operation", "reconcile", "specifies the reconcile operation to trigger")

	genericcli.Must(reconcileCmd.RegisterFlagCompletionFunc("operation", c.Completion.ClusterAdminOperationCompletion))

	// metal admin cluster machine list

	machineListCmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Short:             "list cluster machines",
		ValidArgsFunction: c.Completion.ClusterAdminListCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineList(args)
		},
	}

	// metal admin cluster machine

	machineCmd := &cobra.Command{
		Use:   "machine",
		Short: "commands for cluster machines",
	}

	machineCmd.AddCommand(machineListCmd)

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, reconcileCmd, logsCmd, machineCmd)
}

// TODO: implement firewall ssh, machine/firewall list

func (c *cluster) Create(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv1.ClusterServiceGetRequest{
		Uuid: id,
	}

	resp, err := c.c.Client.Adminv1().Cluster().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) List() ([]*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var labels map[string]string
	if viper.IsSet("labels") {
		tm := tag.NewTagMap(viper.GetStringSlice("labels"))
		labels = map[string]string(tm)
	}

	req := &adminv1.ClusterServiceListRequest{
		Uuid:      pointer.PointerOrNil(viper.GetString("id")),
		Project:   pointer.PointerOrNil(viper.GetString("project")),
		Tenant:    pointer.PointerOrNil(viper.GetString("tenant")),
		Partition: pointer.PointerOrNil(viper.GetString("partition")),
		Seed:      pointer.PointerOrNil(viper.GetString("seed")),
		Name:      pointer.PointerOrNil(viper.GetString("name")),
		Purpose:   pointer.PointerOrNil(viper.GetString("purpose")),
		Labels:    labels,
	}

	resp, err := c.c.Client.Adminv1().Cluster().List(ctx, connect.NewRequest(req))
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
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	expiration := viper.GetDuration("expiration")
	req := &adminv1.ClusterServiceCredentialsRequest{
		Uuid:       id,
		Expiration: durationpb.New(expiration),
	}

	resp, err := c.c.Client.Adminv1().Cluster().Credentials(ctx, connect.NewRequest(req))
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

func (c *cluster) reconcile(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	var operation adminv1.Operate

	switch op := viper.GetString("operation"); op {
	case "reconcile":
		operation = adminv1.Operate_OPERATE_RECONCILE
	case "maintain":
		operation = adminv1.Operate_OPERATE_MAINTAIN
	case "retry":
		operation = adminv1.Operate_OPERATE_RETRY
	default:
		return fmt.Errorf("unsupported operation: %s", op)
	}

	req := &adminv1.ClusterServiceOperateRequest{
		Uuid:    id,
		Operate: operation,
	}

	resp, err := c.c.Client.Adminv1().Cluster().Operate(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to reconcile cluster: %w", err)
	}

	return c.c.DescribePrinter.Print(resp.Msg.Cluster)
}

func (c *cluster) logs(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	req := &adminv1.ClusterServiceGetRequest{
		Uuid: id,
	}

	resp, err := c.c.Client.Adminv1().Cluster().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	return c.c.ListPrinter.Print(resp.Msg.Cluster.Status.LastErrors)
}

func (c *cluster) machineList(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	req := &adminv1.ClusterServiceGetRequest{
		Uuid:         id,
		WithMachines: true,
	}

	resp, err := c.c.Client.Adminv1().Cluster().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	return c.c.ListPrinter.Print(resp.Msg)
}
