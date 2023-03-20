package v1

import (
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type cluster struct {
	c *config.Config
}

func NewClusterCmd(c *config.Config) *cobra.Command {
	w := &cluster{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*connect.Request[apiv1.ClusterServiceCreateRequest], *apiv1.Cluster, *apiv1.Cluster]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*connect.Request[apiv1.ClusterServiceCreateRequest], *apiv1.Cluster, *apiv1.Cluster](w).WithFS(c.Fs),
		Singular:        "cluster",
		Plural:          "clusters",
		Description:     "a cluster of metal-stack cloud",
		Sorter:          sorters.ClusterSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project from where clusters should be listed")
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("name", "", "", "name of the cluster")
			cmd.Flags().StringP("project", "", "", "project where the cluster should be created")
			cmd.Flags().StringP("partition", "", "", "partition of the cluster")
			cmd.Flags().StringP("kubernetes", "", "", "kubernetes version of the cluster")
			cmd.Flags().StringP("workername", "", "", "name of the worker running the cluster")
			cmd.Flags().StringP("machinetype", "", "", "machine type of the worker")
			cmd.Flags().StringP("minsize", "", "", "minimal workers of the cluster")
			cmd.Flags().StringP("maxsize", "", "", "maximal workers of the cluster")
			cmd.Flags().StringP("maintenancebegin", "", "", "time for a possible maintenance begin for the worker")
			cmd.Flags().StringP("maintenanceduration", "", "", "duration of a possible maintenance for the worker")
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project of the cluster")
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "", "", "project of the cluster")
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "uuid of the cluster")
			cmd.Flags().StringP("project", "", "", "project of the cluster")
			cmd.Flags().StringP("kubernetes", "", "", "kubernetes version of the cluster")
			// cmd.Flags().StringP("workername", "", "", "name of the worker running the cluster")
			// cmd.Flags().StringP("machinetype", "", "", "machine type of the worker")
			// cmd.Flags().StringP("minsize", "", "", "minimal workers of the cluster")
			// cmd.Flags().StringP("maxsize", "", "", "maximal workers of the cluster")
			// cmd.Flags().StringP("maxsurge", "", "", "maximal surge of the clusters")
			// cmd.Flags().StringP("maxunavailable", "", "", "maximal unavailable workers of the cluster")
			cmd.Flags().StringP("maintenancebegin", "", "", "time for a possible maintenance begin for the worker")
			cmd.Flags().StringP("maintenanceduration", "", "", "duration of a possible maintenance for the worker")
		},
		CreateRequestFromCLI: func() (*connect.Request[apiv1.ClusterServiceCreateRequest], error) {

			parsedMaintenanceBegin := parseMaintenanceBegin(viper.GetString("maintenancebegin"))

			parsedMaintenanceDuration := parseMaintenanceDuration(viper.GetString("maintenanceduration"))

			clustercr := &apiv1.ClusterServiceCreateRequest{
				Name:      viper.GetString("name"),
				Project:   viper.GetString("project"),
				Partition: viper.GetString("partition"),
				Kubernetes: &apiv1.KubernetesSpec{
					Version: viper.GetString("kubernetes"),
				},
				Workers: []*apiv1.Worker{
					{
						Name:        viper.GetString("workername"),
						MachineType: viper.GetString("machinetype"),
						Minsize:     viper.GetUint32("minsize"),
						Maxsize:     viper.GetUint32("maxsize"),
					},
				},
				Maintenance: &apiv1.Maintenance{
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin:    &parsedMaintenanceBegin,
						Duration: &parsedMaintenanceDuration,
					},
				},
			}
			return connect.NewRequest(clustercr), nil
		},
		UpdateRequestFromCLI: w.updateFromCLI,
	}
	return genericcli.NewCmds(cmdsConfig)
}

func (c *cluster) updateFromCLI(args []string) (*apiv1.Cluster, error) {
	clusterToUpdate, err := c.Get(viper.GetString("uuid"))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve cluster: %w", err)
	}
	if viper.IsSet("project") {
		clusterToUpdate.Project = viper.GetString("project")
	}
	if clusterToUpdate.Project == "" {
		return nil, fmt.Errorf("project must be provided")
	}
	if viper.IsSet("kubernetes") {
		clusterToUpdate.Kubernetes.Version = viper.GetString("kubernetes")
	}
	//TODO Should it be possible to update workers too?
	// if viper.IsSet("workername") {
	// 	clusterToUpdate.GetWorkers()
	// }
	if viper.IsSet("maintenancebegin") {
		parsedMaintenance := parseMaintenanceBegin(viper.GetString("maintenancebegin"))
		clusterToUpdate.Maintenance.TimeWindow.Begin = &parsedMaintenance
		fmt.Println(clusterToUpdate)
	}
	if viper.IsSet("maintenanceduration") {
		parsedDuration := parseMaintenanceDuration(viper.GetString("maintenanceduration"))
		clusterToUpdate.Maintenance.TimeWindow.Duration = &parsedDuration
	}
	return clusterToUpdate, nil
}

func (c *cluster) Create(rq *connect.Request[apiv1.ClusterServiceCreateRequest]) (*apiv1.Cluster, error) {
	resp, err := c.c.Apiv1Client.Cluster().Create(c.c.Ctx, rq)
	if err != nil {
		return nil, err
	}
	return resp.Msg.Cluster, nil
}

func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	project := viper.GetString("project")
	if project == "" {
		return nil, fmt.Errorf("project must be provided")
	}
	resp, err := c.c.Apiv1Client.Cluster().Delete(c.c.Ctx, connect.NewRequest(&apiv1.ClusterServiceDeleteRequest{
		Uuid:    id,
		Project: project,
	}))

	if err != nil {
		return nil, err
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	project := viper.GetString("project")
	if project == "" {
		return nil, fmt.Errorf("project must be provided")
	}
	resp, err := c.c.Apiv1Client.Cluster().Get(c.c.Ctx, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
		Uuid:    id,
		Project: project,
	}))

	if err != nil {
		return nil, err
	}

	return resp.Msg.Cluster, nil
}

func (c *cluster) List() ([]*apiv1.Cluster, error) {
	project := viper.GetString("project")
	if project == "" {
		return nil, fmt.Errorf("project must be provided")
	}

	resp, err := c.c.Apiv1Client.Cluster().List(c.c.Ctx, connect.NewRequest(&apiv1.ClusterServiceListRequest{
		Project: project,
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Clusters, nil
}

// TODO implement update command and write tests
func (c *cluster) Update(rq *apiv1.Cluster) (*apiv1.Cluster, error) {
	// var workerUpdate []*apiv1.WorkerUpdate
	// for _, worker := range rq.Workers {
	// 	workerUpdate = append(workerUpdate, &apiv1.WorkerUpdate{
	// 		Name:           worker.Name,
	// 		MachineType:    &worker.MachineType,
	// 		Minsize:        &worker.Minsize,
	// 		Maxsize:        &worker.Maxsize,
	// 		Maxsurge:       &worker.Maxsurge,
	// 		Maxunavailable: &worker.Maxunavailable,
	// 	})
	// }
	// TODO I guess the update of the maintenance field is not working right now.
	resp, err := c.c.Apiv1Client.Cluster().Update(c.c.Ctx, &connect.Request[apiv1.ClusterServiceUpdateRequest]{
		Msg: &apiv1.ClusterServiceUpdateRequest{
			Uuid:       rq.Uuid,
			Project:    rq.Project,
			Kubernetes: rq.Kubernetes,
			Maintenance: &apiv1.Maintenance{
				TimeWindow: &apiv1.MaintenanceTimeWindow{
					Begin: &timestamppb.Timestamp{
						Seconds: rq.Maintenance.TimeWindow.Begin.Seconds,
					},
					Duration: &durationpb.Duration{
						Seconds: rq.Maintenance.TimeWindow.Duration.Seconds,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Msg.Cluster, nil
}

func (c *cluster) ToCreate(r *apiv1.Cluster) (*connect.Request[apiv1.ClusterServiceCreateRequest], error) {
	return clusterResponseToCreate(r), nil
}

func (c *cluster) ToUpdate(r *apiv1.Cluster) (*apiv1.Cluster, error) {
	return clusterResponseToUpdate(r), nil
}

func clusterResponseToCreate(r *apiv1.Cluster) *connect.Request[apiv1.ClusterServiceCreateRequest] {
	return &connect.Request[apiv1.ClusterServiceCreateRequest]{
		Msg: &apiv1.ClusterServiceCreateRequest{
			Name:        r.Name,
			Project:     r.Project,
			Partition:   r.Partition,
			Kubernetes:  r.Kubernetes,
			Workers:     r.Workers,
			Maintenance: r.Maintenance,
		},
	}
}

func clusterResponseToUpdate(r *apiv1.Cluster) *apiv1.Cluster {
	return r
}

func parseMaintenanceBegin(m string) timestamppb.Timestamp {
	const format = "Jan 2, 2006 at 3:04pm (MST)"
	t, _ := time.Parse(format, "Jan 1, 1970 at "+m+" (PST)")
	timeInUnix := t.Unix()
	timeInUTC := time.Unix(timeInUnix, 0).UTC()
	ts := timestamppb.New(timeInUTC)
	return *ts
}

func parseMaintenanceDuration(m string) durationpb.Duration {
	hours, _ := time.ParseDuration(m)
	duration := durationpb.New(time.Duration(hours))
	return *duration
}
