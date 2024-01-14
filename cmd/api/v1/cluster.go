package v1

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/fatih/color"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/kubernetes"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack-cloud/cli/pkg/helpers"
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

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.ClusterServiceCreateRequest, *apiv1.ClusterServiceUpdateRequest, *apiv1.Cluster]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*apiv1.ClusterServiceCreateRequest, *apiv1.ClusterServiceUpdateRequest, *apiv1.Cluster](w).WithFS(c.Fs),
		Singular:        "cluster",
		Plural:          "clusters",
		Description:     "manage kubernetes clusters",
		Sorter:          sorters.ClusterSorter(),
		ValidArgsFn:     c.Completion.ClusterListCompletion,
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the clusters")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		CreateRequestFromCLI: w.createFromCLI,
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().String("name", "", "name of the cluster")
			cmd.Flags().StringP("project", "p", "", "project of the cluster")
			cmd.Flags().String("partition", "", "partition of the cluster")
			cmd.Flags().String("kubernetes-version", "", "kubernetes version of the cluster")
			cmd.Flags().Int32("maintenance-hour", 0, "hour in which cluster maintenance is allowed to take place")
			cmd.Flags().Int32("maintenance-minute", 0, "minute in which cluster maintenance is allowed to take place")
			cmd.Flags().String("maintenance-timezone", time.Local.String(), "timezone used for the maintenance time window") // nolint
			cmd.Flags().Duration("maintenance-duration", 2*time.Hour, "duration in which cluster maintenance is allowed to take place")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("kubernetes-version", c.Completion.KubernetesVersionAssetListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the cluster")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the cluster")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the cluster")
			cmd.Flags().String("kubernetes-version", "", "kubernetes version of the cluster")
			cmd.Flags().Uint32("maintenance-hour", 0, "hour in which cluster maintenance is allowed to take place")
			cmd.Flags().Uint32("maintenance-minute", 0, "minute in which cluster maintenance is allowed to take place")
			cmd.Flags().String("maintenance-timezone", time.Local.String(), "timezone used for the maintenance time window") // nolint
			cmd.Flags().Duration("maintenance-duration", 2*time.Hour, "duration in which cluster maintenance is allowed to take place")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("kubernetes-version", c.Completion.KubernetesVersionAssetListCompletion))
		},
		UpdateRequestFromCLI: w.updateFromCLI,
	}

	kubeconfigCmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "fetch kubeconfig of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.kubeconfig(args)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	kubeconfigCmd.Flags().StringP("project", "p", "", "the project in which the cluster resides for which to get the kubeconfig for")
	kubeconfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")
	kubeconfigCmd.Flags().Bool("merge", true, "merges the kubeconfig into default kubeconfig instead of printing it to the console")
	kubeconfigCmd.Flags().String("kubeconfig", "", "specify an explicit path for the merged kubeconfig to be written, defaults to default kubeconfig paths if not provided")

	genericcli.Must(kubeconfigCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

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

	monitoringCmd.Flags().String("project", "", "the project in which the cluster resides for which to get the kubeconfig for")

	genericcli.Must(monitoringCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, monitoringCmd)
}

func (c *cluster) Create(req *apiv1.ClusterServiceCreateRequest) (*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	if req.Partition == "" {
		return nil, fmt.Errorf("partition is required")
	}

	resp, err := c.c.Client.Apiv1().Cluster().Create(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) createFromCLI() (*apiv1.ClusterServiceCreateRequest, error) {
	rq := &apiv1.ClusterServiceCreateRequest{
		Name:        viper.GetString("name"),
		Project:     c.c.GetProject(),
		Partition:   viper.GetString("partition"),
		Maintenance: &apiv1.Maintenance{},
	}

	if viper.IsSet("kubernetes-version") {
		rq.Kubernetes = &apiv1.KubernetesSpec{
			Version: viper.GetString("kubernetes-version"),
		}
	}

	if viper.IsSet("maintenance-hour") {
		rq.Maintenance.TimeWindow = &apiv1.MaintenanceTimeWindow{
			Begin: &apiv1.Time{
				Hour:     viper.GetUint32("maintenance-hour"),
				Minute:   viper.GetUint32("maintenance-minute"),
				Timezone: viper.GetString("maintenance-timezone"),
			},
			Duration: durationpb.New(viper.GetDuration("maintenance-duration")),
		}
	}

	return rq, nil
}

func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.ClusterServiceDeleteRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	if viper.IsSet("file") {
		var err error
		req.Uuid, req.Project, err = helpers.DecodeProject(id)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Apiv1().Cluster().Delete(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to delete cluster: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.ClusterServiceGetRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Cluster().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) List() ([]*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.ClusterServiceListRequest{
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Cluster().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	return resp.Msg.Clusters, nil
}

func (c *cluster) Convert(r *apiv1.Cluster) (string, *apiv1.ClusterServiceCreateRequest, *apiv1.ClusterServiceUpdateRequest, error) {
	return helpers.EncodeProject(r.Uuid, r.Project), ClusterResponseToCreate(r), ClusterResponseToUpdate(r), nil
}

func ClusterResponseToCreate(r *apiv1.Cluster) *apiv1.ClusterServiceCreateRequest {
	return &apiv1.ClusterServiceCreateRequest{
		Name:        r.Name,
		Project:     r.Project,
		Partition:   r.Partition,
		Kubernetes:  r.Kubernetes,
		Workers:     r.Workers,
		Maintenance: r.Maintenance,
	}
}

func ClusterResponseToUpdate(r *apiv1.Cluster) *apiv1.ClusterServiceUpdateRequest {
	return &apiv1.ClusterServiceUpdateRequest{
		Uuid:       r.Uuid,
		Project:    r.Project,
		Kubernetes: r.Kubernetes,
		// Workers:     workers, // TODO
		Maintenance: r.Maintenance,
	}
}

func (c *cluster) Update(req *apiv1.ClusterServiceUpdateRequest) (*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().Cluster().Update(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to update cluster: %w", err)
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) updateFromCLI(args []string) (*apiv1.ClusterServiceUpdateRequest, error) {
	uuid, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	cluster, err := c.Get(uuid)
	if err != nil {
		return nil, err
	}

	rq := &apiv1.ClusterServiceUpdateRequest{
		Uuid:        uuid,
		Project:     c.c.GetProject(),
		Kubernetes:  &apiv1.KubernetesSpec{},
		Maintenance: &apiv1.Maintenance{},
	}

	if viper.IsSet("maintenance-hour") || viper.IsSet("maintenance-minute") || viper.IsSet("maintenance-duration") {
		rq.Maintenance = cluster.Maintenance

		if viper.IsSet("maintenance-hour") {
			rq.Maintenance.TimeWindow.Begin.Hour = viper.GetUint32("maintenance-hour")
			rq.Maintenance.TimeWindow.Begin.Timezone = viper.GetString("maintenance-timezone")

		}
		if viper.IsSet("maintenance-minute") {
			rq.Maintenance.TimeWindow.Begin.Minute = viper.GetUint32("maintenance-minute")
			rq.Maintenance.TimeWindow.Begin.Timezone = viper.GetString("maintenance-timezone")
		}
		if viper.IsSet("maintenance-duration") {
			rq.Maintenance.TimeWindow.Duration = durationpb.New(viper.GetDuration("maintenance-duration"))
		}
	}

	if viper.IsSet("kubernetes-version") {
		rq.Kubernetes.Version = viper.GetString("kubernetes-version")
	}

	return rq, nil
}

func (c *cluster) kubeconfig(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	req := &apiv1.ClusterServiceGetCredentialsRequest{
		Uuid:       id,
		Project:    c.c.GetProject(),
		Expiration: durationpb.New(viper.GetDuration("expiration")),
	}

	resp, err := c.c.Client.Apiv1().Cluster().GetCredentials(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to get cluster credentials: %w", err)
	}

	if !viper.GetBool("merge") {
		fmt.Fprintln(c.c.Out, resp.Msg.Kubeconfig)
		return nil
	}

	projectResp, err := c.c.Client.Apiv1().Project().Get(ctx, connect.NewRequest(&apiv1.ProjectServiceGetRequest{Project: c.c.GetProject()}))
	if err != nil {
		return err
	}

	var (
		kubeconfigPath = viper.GetString("kubeconfig")
		projectName    = helpers.TrimProvider(projectResp.Msg.Project.Name)
	)

	merged, err := kubernetes.MergeKubeconfig(c.c.Fs, []byte(resp.Msg.Kubeconfig), pointer.PointerOrNil(kubeconfigPath), &projectName)
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
