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
	"github.com/metal-stack-cloud/cli/pkg/helpers"
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
		ValidArgsFn:     c.Completion.AdminClusterListCompletion,
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

			genericcli.Must(cmd.RegisterFlagCompletionFunc("id", c.Completion.AdminClusterListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("name", c.Completion.AdminClusterNameListCompletion))
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
		ValidArgsFunction: c.Completion.AdminClusterListCompletion,
	}

	kubeconfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")
	kubeconfigCmd.Flags().Bool("merge", true, "merges the kubeconfig into default kubeconfig instead of printing it to the console")
	kubeconfigCmd.Flags().String("kubeconfig", "", "specify an explicit path for the merged kubeconfig to be written, defaults to default kubeconfig paths if not provided")

	// metal admin cluster machine list

	machineListCmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Short:             "list cluster machines",
		ValidArgsFunction: c.Completion.AdminClusterListCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineList(args)
		},
	}

	// metal admin cluster machine ssh

	machineSSHCmd := &cobra.Command{
		Use:               "ssh",
		Short:             "ssh to cluster machines",
		ValidArgsFunction: c.Completion.AdminClusterListCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.machineSSH(args)
		},
	}

	machineSSHCmd.Flags().String("machine-id", "", "the firewall's machine id to connect to")

	genericcli.Must(machineSSHCmd.RegisterFlagCompletionFunc("machine-id", c.Completion.AdminClusterFirewallListCompletion))

	// metal admin cluster machine

	machineCmd := &cobra.Command{
		Use:   "machine",
		Short: "commands for cluster machines",
	}

	machineCmd.AddCommand(machineListCmd, machineSSHCmd)

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, machineCmd)
}

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
		_, _ = fmt.Fprintln(c.c.Out, resp.Msg.Kubeconfig)
		return nil
	}

	var (
		kubeconfigPath = viper.GetString("kubeconfig")
	)

	merged, err := kubernetes.MergeKubeconfig(c.c.Fs, []byte(resp.Msg.Kubeconfig), pointer.PointerOrNil(kubeconfigPath), nil, c.c.GetProject(), id) // FIXME: reverse lookup project name
	if err != nil {
		return err
	}

	err = afero.WriteFile(c.c.Fs, merged.Path, merged.Raw, 0600)
	if err != nil {
		return fmt.Errorf("unable to write merged kubeconfig: %w", err)
	}

	_, _ = fmt.Fprintf(c.c.Out, "%s merged context %q into %s\n", color.GreenString("✔"), merged.ContextName, merged.Path)

	return nil
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

func (c *cluster) machineSSH(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	var (
		machineID = viper.GetString("machine-id")
	)

	if machineID == "" {
		return fmt.Errorf("machine id is required")
	}

	credsResp, err := c.c.Client.Adminv1().Cluster().Credentials(ctx, connect.NewRequest(&adminv1.ClusterServiceCredentialsRequest{
		Uuid:    id,
		WithSsh: true,
		WithVpn: true,
	}))
	if err != nil {
		return fmt.Errorf("failed to get cluster credentials: %w", err)
	}

	return helpers.SSHViaVPN(c.c.Out, machineID, credsResp.Msg)
}
