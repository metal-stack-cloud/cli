package tableprinters

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) VolumeTable(data []*apiv1.Volume, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Name", "Size", "Usage", "Replicas", "StorageClass", "Project", "Partition"}
	if wide {
		header = []string{"ID", "Name", "Size", "Usage", "Replicas", "StorageClass", "Project", "Partition", "Nodes"}
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

		short := []string{volumeID, name, size, usage, replica, sc, project, partition}
		if wide {
			short := append(short, strings.Join(nodes, "\n"))

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
		if len(parts) >= 1 {
			node := strings.TrimSuffix(parts[1], ".node")
			nodes = append(nodes, node)
		}
	}
	return nodes
}
