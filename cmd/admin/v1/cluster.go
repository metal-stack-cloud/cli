package v1

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Description:     "cluster related actions of metalstack.cloud",
		Sorter:          sorters.ClusterSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		OnlyCmds:        genericcli.OnlyCmds(genericcli.DescribeCmd, genericcli.ListCmd),

		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().BoolP("machines", "", false, "show machines of a cluster")
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
			req := &adminv1.ClusterServiceCredentialsRequest{
				Uuid: id,
			}
			resp, err := c.Adminv1Client.Cluster().Credentials(c.Ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to get cluster credentials: %w", err)
			}

			fmt.Println(resp.Msg.Kubeconfig)
			return nil
		},
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
			resp, err := c.Adminv1Client.Cluster().Operate(c.Ctx, connect.NewRequest(req))
			if err != nil {
				return fmt.Errorf("failed to reconcile cluster: %w", err)
			}

			return c.ListPrinter.Print(resp.Msg.Cluster)
		},
	}

	reconcileCmd.Flags().Bool("reconcile", true, "trigger cluster reconciliation")
	reconcileCmd.Flags().Bool("maintain", false, "trigger cluster maintain reconciliation")
	reconcileCmd.Flags().Bool("retry", false, "trigger cluster retry reconciliation")

	return genericcli.NewCmds(cmdsConfig, kubeconfigCmd, reconcileCmd)
}

// TODO: implement GetCredentials, Logs, Operate, firewall ssh, machine/firewall list

// Create implements genericcli.CRUD
func (c *cluster) Create(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

// Delete implements genericcli.CRUD
func (c *cluster) Delete(id string) (*apiv1.Cluster, error) {
	panic("unimplemented")
}

// Get implements genericcli.CRUD
func (c *cluster) Get(id string) (*apiv1.Cluster, error) {
	showMachines := viper.GetBool("machines")
	req := &adminv1.ClusterServiceGetRequest{
		Uuid:         id,
		WithMachines: showMachines,
	}
	resp, err := c.c.Adminv1Client.Cluster().Get(c.c.Ctx, connect.NewRequest(req))
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

// List implements genericcli.CRUD
func (c *cluster) List() ([]*apiv1.Cluster, error) {
	// FIXME implement filters and paging

	req := &adminv1.ClusterServiceListRequest{}
	resp, err := c.c.Adminv1Client.Cluster().List(c.c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}
	return resp.Msg.Clusters, nil
}

// ToCreate implements genericcli.CRUD
func (c *cluster) ToCreate(r *apiv1.Cluster) (any, error) {
	panic("unimplemented")
}

// ToUpdate implements genericcli.CRUD
func (c *cluster) ToUpdate(r *apiv1.Cluster) (any, error) {
	panic("unimplemented")
}

// Update implements genericcli.CRUD
func (c *cluster) Update(rq any) (*apiv1.Cluster, error) {
	panic("unimplemented")
}
