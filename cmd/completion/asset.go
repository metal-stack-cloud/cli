package completion

import (
	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) PartitionAssetListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.AssetServiceListRequest{}
	resp, err := c.Client.Apiv1().Asset().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var names []string
	for _, asset := range resp.Msg.Assets {
		asset := asset
		for partition := range asset.Region.Partitions {
			names = append(names, partition)
		}
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) KubernetesVersionAssetListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.AssetServiceListRequest{}
	resp, err := c.Client.Apiv1().Asset().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var versions []string
	for _, asset := range resp.Msg.Assets {
		asset := asset
		for _, kubernetes := range asset.Kubernetes {
			kubernetes := kubernetes
			versions = append(versions, kubernetes.Version)
		}
	}
	return versions, cobra.ShellCompDirectiveNoFileComp
}

func (c *Completion) MachineTypeAssetListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &apiv1.AssetServiceListRequest{}

	resp, err := c.Client.Apiv1().Asset().List(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var relevantRegions []*apiv1.Asset
	for _, asset := range resp.Msg.Assets {
		asset := asset

		if partition := cmd.Flag("partition").Value.String(); partition != "" {
			_, ok := asset.Region.Partitions[partition]
			if !ok {
				continue
			}
		}

		relevantRegions = append(relevantRegions, asset)
	}

	var types []string
	for _, region := range relevantRegions {
		region := region

		for _, machineType := range region.MachineTypes {
			machineType := machineType
			types = append(types, machineType.Name)
		}
	}

	return types, cobra.ShellCompDirectiveNoFileComp
}
