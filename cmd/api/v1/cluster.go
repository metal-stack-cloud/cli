package v1

import (
	"encoding/json"
	"fmt"
	"slices"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			cmd.Flags().String("worker-group", "group-0", "the name of the initial worker group")
			cmd.Flags().Uint32("worker-min", 1, "the minimum amount of worker nodes of the worker group")
			cmd.Flags().Uint32("worker-max", 3, "the maximum amount of worker nodes of the worker group")
			cmd.Flags().Uint32("worker-max-surge", 1, "the maximum amount of new worker nodes added to the worker group during a rolling update")
			cmd.Flags().Uint32("worker-max-unavailable", 0, "the maximum amount of worker nodes removed from the worker group during a rolling update")
			cmd.Flags().String("worker-type", "", "the worker type of the initial worker group")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("kubernetes-version", c.Completion.KubernetesVersionAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("worker-type", c.Completion.MachineTypeAssetListCompletion))
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
			cmd.Flags().String("worker-group", "", "the name of the worker group to add, update or remove")
			cmd.Flags().Uint32("worker-min", 1, "the minimum amount of worker nodes of the worker group")
			cmd.Flags().Uint32("worker-max", 3, "the maximum amount of worker nodes of the worker group")
			cmd.Flags().Uint32("worker-max-surge", 1, "the maximum amount of new worker nodes added to the worker group during a rolling update")
			cmd.Flags().Uint32("worker-max-unavailable", 0, "the maximum amount of worker nodes removed from the worker group during a rolling update")
			cmd.Flags().String("worker-type", "", "the worker type of the initial worker group")
			cmd.Flags().Bool("remove-worker-group", false, "if set the selected worker group is being removed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("kubernetes-version", c.Completion.KubernetesVersionAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("worker-type", c.Completion.MachineTypeAssetListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("worker-group", c.Completion.ClusterWorkerGroupsCompletion))
		},
		UpdateRequestFromCLI: w.updateFromCLI,
	}

	// cluster kubeconfig

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

	execConfigCmd := &cobra.Command{
		Use:   "exec-config",
		Short: "fetch exec-config of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.execConfig(args)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	execConfigCmd.Flags().StringP("project", "p", "", "the project in which the cluster resides for which to get the kubeconfig for")
	execConfigCmd.Flags().DurationP("expiration", "", 8*time.Hour, "kubeconfig will expire after given time")

	genericcli.Must(execConfigCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

	// cluster monitoring

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

	// cluster reconcile

	reconcileCmd := &cobra.Command{
		Use:   "reconcile",
		Short: "reconcile a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.reconcile(args)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	reconcileCmd.Flags().String("operation", "reconcile", "specifies the reconcile operation to trigger")
	reconcileCmd.Flags().StringP("project", "p", "", "project of the cluster")

	genericcli.Must(reconcileCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
	genericcli.Must(reconcileCmd.RegisterFlagCompletionFunc("operation", c.Completion.ClusterOperationCompletion))

	// metal cluster status

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "fetch status of a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.status(args)
		},
		ValidArgsFunction: c.Completion.ClusterListCompletion,
	}

	statusCmd.Flags().StringP("project", "p", "", "project of the cluster")

	genericcli.Must(statusCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, execConfigCmd, monitoringCmd, reconcileCmd, statusCmd)
}

func (c *cluster) Create(req *apiv1.ClusterServiceCreateRequest) (*apiv1.Cluster, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	if req.Partition == "" {
		return nil, fmt.Errorf("partition is required")
	}

	resp, err := c.c.Client.Apiv1().Cluster().Create(ctx, connect.NewRequest(req))
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			return nil, genericcli.AlreadyExistsError()
		}
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

	if helpers.IsAnyViperFlagSet("worker-group", "worker-min", "worker-max", "worker-max-surge", "worker-max-unavailable", "worker-type") {
		rq.Workers = append(rq.Workers, &apiv1.Worker{
			Name:           viper.GetString("worker-group"),
			MachineType:    viper.GetString("worker-type"),
			Minsize:        viper.GetUint32("worker-min"),
			Maxsize:        viper.GetUint32("worker-max"),
			Maxsurge:       viper.GetUint32("worker-max-surge"),
			Maxunavailable: viper.GetUint32("worker-max-unavailable"),
		})
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

	if !viper.GetBool("skip-security-prompts") {
		cluster, err := c.Get(id)
		if err != nil {
			return nil, err
		}

		if err := genericcli.PromptCustom(&genericcli.PromptConfig{
			Message:         fmt.Sprintf(`Do you really want to delete "%s"? This operation cannot be undone.`, color.RedString(cluster.Name)),
			ShowAnswers:     true,
			AcceptedAnswers: []string{"y", "yes"},
			DefaultAnswer:   "n",
			No:              "n",
			In:              c.c.In,
			Out:             c.c.Out,
		}); err != nil {
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
		Uuid:        r.Uuid,
		Project:     r.Project,
		Kubernetes:  r.Kubernetes,
		Workers:     clusterWorkersToWorkerUpdate(r.Workers),
		Maintenance: r.Maintenance,
	}
}

func clusterWorkersToWorkerUpdate(workers []*apiv1.Worker) []*apiv1.WorkerUpdate {
	var res []*apiv1.WorkerUpdate
	for _, worker := range workers {
		worker := worker

		res = append(res, clusterWorkerToWorkerUpdate(worker))
	}

	return res
}

func clusterWorkerToWorkerUpdate(worker *apiv1.Worker) *apiv1.WorkerUpdate {
	return &apiv1.WorkerUpdate{
		Name:           worker.Name,
		MachineType:    pointer.Pointer(worker.MachineType),
		Minsize:        pointer.Pointer(worker.Minsize),
		Maxsize:        pointer.Pointer(worker.Maxsize),
		Maxsurge:       pointer.Pointer(worker.Maxsurge),
		Maxunavailable: pointer.Pointer(worker.Maxunavailable),
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
		Uuid:    uuid,
		Project: cluster.Project,
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
		rq.Kubernetes = cluster.Kubernetes

		rq.Kubernetes.Version = viper.GetString("kubernetes-version")
	}

	findWorkerGroup := func() (*apiv1.Worker, error) {
		if viper.GetString("worker-group") == "" {
			if len(cluster.Workers) != 1 {
				return nil, fmt.Errorf("please specify the group to act on using the flag --worker-group")
			}

			return cluster.Workers[0], nil
		}

		for _, worker := range cluster.Workers {
			worker := worker
			if worker.Name == viper.GetString("worker-group") {
				return worker, nil
			}
		}

		return nil, nil
	}

	if helpers.IsAnyViperFlagSet("worker-group", "worker-min", "worker-max", "worker-max-surge", "worker-max-unavailable", "worker-type", "remove-worker-group") {
		type operation string

		const (
			update operation = "Updating"
			delete operation = "Deleting"
			add    operation = "Adding"
		)

		var (
			newWorkers []*apiv1.WorkerUpdate
			showPrompt = func(op operation, name string) error {
				if viper.GetBool("skip-security-prompts") {
					return nil
				}

				return genericcli.PromptCustom(&genericcli.PromptConfig{
					Message:     fmt.Sprintf("%s worker group %q, continue?", op, name),
					ShowAnswers: true,
					Out:         c.c.PromptOut,
					In:          c.c.In,
				})
			}
		)

		selectedGroup, err := findWorkerGroup()
		if err != nil {
			return nil, err
		}

		if selectedGroup == nil {
			if viper.IsSet("remove-worker-group") {
				return nil, fmt.Errorf("cluster has no worker group with name %q", viper.GetString("worker-group"))
			}

			if err := showPrompt(add, viper.GetString("worker-group")); err != nil {
				return nil, err
			}

			newWorkers = append(clusterWorkersToWorkerUpdate(cluster.Workers), &apiv1.WorkerUpdate{
				Name:           viper.GetString("worker-group"),
				MachineType:    pointer.PointerOrNil(viper.GetString("worker-type")),
				Minsize:        pointer.PointerOrNil(viper.GetUint32("worker-min")),
				Maxsize:        pointer.PointerOrNil(viper.GetUint32("worker-max")),
				Maxsurge:       pointer.PointerOrNil(viper.GetUint32("worker-max-surge")),
				Maxunavailable: pointer.PointerOrNil(viper.GetUint32("worker-max-unavailable")),
			})
		} else {
			if viper.IsSet("remove-worker-group") {
				if err := showPrompt(delete, selectedGroup.Name); err != nil {
					return nil, err
				}

				newWorkers = clusterWorkersToWorkerUpdate(cluster.Workers)
				newWorkers = slices.DeleteFunc(newWorkers, func(w *apiv1.WorkerUpdate) bool {
					return w.Name == selectedGroup.Name
				})
			} else {
				if err := showPrompt(update, selectedGroup.Name); err != nil {
					return nil, err
				}

				for _, worker := range cluster.Workers {
					worker := worker

					workerUpdate := clusterWorkerToWorkerUpdate(worker)

					if worker.Name == selectedGroup.Name {
						if viper.IsSet("worker-min") {
							workerUpdate.Minsize = pointer.Pointer(viper.GetUint32("worker-min"))
						}
						if viper.IsSet("worker-max") {
							workerUpdate.Maxsize = pointer.Pointer(viper.GetUint32("worker-max"))
						}
						if viper.IsSet("worker-max-surge") {
							workerUpdate.Maxsurge = pointer.Pointer(viper.GetUint32("worker-max-surge"))
						}
						if viper.IsSet("worker-max-unavailable") {
							workerUpdate.Maxunavailable = pointer.Pointer(viper.GetUint32("worker-max-unavailable"))
						}
						if viper.IsSet("worker-type") {
							workerUpdate.MachineType = pointer.Pointer(viper.GetString("worker-type"))
						}
					}

					newWorkers = append(newWorkers, workerUpdate)
				}
			}
		}

		rq.Workers = newWorkers
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
		_, _ = fmt.Fprintln(c.c.Out, resp.Msg.Kubeconfig)
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

	merged, err := kubernetes.MergeKubeconfig(c.c.Fs, []byte(resp.Msg.Kubeconfig), pointer.PointerOrNil(kubeconfigPath), &projectName, projectResp.Msg.Project.Uuid, id)
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

func (c *cluster) execConfig(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ec, err := kubernetes.NewUserExecCache(c.c.Fs)
	if err != nil {
		return err
	}

	creds, err := ec.LoadCachedCredentials(id)
	if err != nil {
		// we cannot load cache, so cleanup the cache
		_ = ec.Clean(id)
	}
	if creds == nil {
		req := &apiv1.ClusterServiceGetCredentialsRequest{
			Uuid:       id,
			Project:    c.c.GetProject(),
			Expiration: durationpb.New(viper.GetDuration("expiration")),
		}

		resp, err := c.c.Client.Apiv1().Cluster().GetCredentials(ctx, connect.NewRequest(req))
		if err != nil {
			return fmt.Errorf("failed to get cluster credentials: %w", err)
		}

		// the kubectl client will re-request credentials when the old credentials expire, so
		// the user won't realize if the expiration is short.
		creds, err = ec.ExecConfig(id, resp.Msg.GetKubeconfig(), viper.GetDuration("expiration"))
		if err != nil {
			return fmt.Errorf("unable to decode kubeconfig: %w", err)
		}
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal exec cred: %w", err)
	}
	_, _ = fmt.Fprintf(c.c.Out, "%s\n", data)
	return nil
}

func (c *cluster) reconcile(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	var operation apiv1.Operate

	switch op := viper.GetString("operation"); op {
	case "reconcile":
		operation = apiv1.Operate_OPERATE_RECONCILE
	case "maintain":
		operation = apiv1.Operate_OPERATE_MAINTAIN
	case "retry":
		operation = apiv1.Operate_OPERATE_RETRY
	default:
		return fmt.Errorf("unsupported operation: %s", op)
	}

	req := &apiv1.ClusterServiceOperateRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
		Operate: operation,
	}

	resp, err := c.c.Client.Apiv1().Cluster().Operate(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to reconcile cluster: %w", err)
	}

	return c.c.DescribePrinter.Print(resp.Msg.Cluster)
}

func (c *cluster) status(args []string) error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	req := &apiv1.ClusterServiceGetRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Cluster().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	err = c.c.ListPrinter.Print(resp.Msg.Cluster.Status.Conditions)
	if err != nil {
		return err
	}

	if len(resp.Msg.Cluster.Status.LastErrors) == 0 {
		return nil
	}

	_, _ = fmt.Fprintln(c.c.Out)
	_, _ = fmt.Fprintln(c.c.Out, "Last Errors:")

	return c.c.ListPrinter.Print(resp.Msg.Cluster.Status.LastErrors)
}
