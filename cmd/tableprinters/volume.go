package tableprinters

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) VolumeTable(data []*apiv1.Volume, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Name", "Size", "Usage", "Replicas", "ClusterName", "StorageClass", "Project", "Partition"}
	if wide {
		header = []string{"ID", "Name", "Size", "Usage", "Replicas", "ClusterName", "StorageClass", "Project", "Partition", "Nodes", "Labels"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Uuid < data[j].Uuid })
	for _, vol := range data {
		volumeID := vol.Uuid
		name := vol.Name
		size := humanize.IBytes(vol.Size)
		usage := humanize.IBytes(vol.Usage)
		replica := fmt.Sprintf("%d", vol.ReplicaCount)
		sc := vol.StorageClass
		partition := vol.Partition
		project := vol.Project
		nodes := connectedHosts(vol)
		labels := volumeLabels(vol)
		clusterName := "-"
		if vol.ClusterName != "" {
			clusterName = vol.ClusterName
		}

		short := []string{volumeID, name, size, usage, replica, clusterName, sc, project, partition}
		if wide {
			short := append(short, strings.Join(nodes, "\n"), strings.Join(labels, "\n"))

			rows = append(rows, short)
		} else {
			rows = append(rows, short)
		}
	}

	return header, rows, nil
}

// connectedHosts returns the worker nodes without internal prefixes and suffixes
func connectedHosts(vol *apiv1.Volume) []string {
	nodes := []string{}
	for _, n := range vol.AttachedTo {
		// nqn.2019-09.com.lightbitslabs:host:shoot--pddhz9--duros-tst9-group-0-6b7bb-2cnvs.node
		parts := strings.Split(n, ":host:")
		if len(parts) > 1 {
			node := strings.TrimSuffix(parts[1], ".node")
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func volumeLabels(vol *apiv1.Volume) []string {
	labels := make([]string, 0, len(vol.Labels))

	for _, l := range vol.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", l.Key, l.Value))
	}

	slices.Sort(labels)
	return labels
}
