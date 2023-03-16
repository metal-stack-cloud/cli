package tableprinters

import (
	"sort"

	"github.com/dustin/go-humanize"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) SnapshotTable(data []*apiv1.Snapshot, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Name", "Size", "Usage", "SourceVolumeID", "SourceVolumeName", "Project", "Partition"}

	sort.SliceStable(data, func(i, j int) bool { return data[i].Uuid < data[j].Uuid })
	for _, snap := range data {
		snapshotID := snap.Uuid
		name := snap.Name
		size := humanize.IBytes(snap.Size)
		usage := humanize.IBytes(snap.Usage)
		sourceVolumeID := snap.SourceVolumeUuid
		sourceVolumeName := snap.SourceVolumeName
		partition := snap.Partition
		project := snap.Project

		short := []string{snapshotID, name, size, usage, sourceVolumeID, sourceVolumeName, project, partition}
		rows = append(rows, short)
	}

	return header, rows, nil
}
