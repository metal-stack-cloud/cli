package tableprinters

import (
	"fmt"

	"github.com/dustin/go-humanize"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
)

func (t *TablePrinter) StorageClusterInfoTable(data []*adminv1.StorageClusterInfo, _ bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Partition", "Version", "Health", "Nodes NA", "Volumes D/NA/RO", "Physical Installed/Managed", "Physical Effective/Free/Used", "Logical Total/Used", "Estimated Total/Free", "Compression"}
	)

	for _, info := range data {
		if info == nil || info.Statistics == nil {
			continue
		}

		partition := info.Partition
		health := info.Health.State
		numdegradedvolumes := info.Health.NumDegradedVolumes
		numnotavailablevolumes := info.Health.NumNotAvailableVolumes
		numreadonlyvolumes := info.Health.NumReadOnlyVolumes
		numinactivenodes := info.Health.NumInactiveNodes

		compressionratio := ""
		if info.Statistics != nil {
			ratio := info.Statistics.CompressionRatio
			compressionratio = fmt.Sprintf("%d%%", int(100.0*(1-ratio)))
		}
		effectivephysicalstorage := humanize.IBytes(info.Statistics.EffectivePhysicalStorage)
		freephysicalstorage := humanize.IBytes(info.Statistics.FreePhysicalStorage)
		physicalusedstorage := humanize.IBytes(info.Statistics.PhysicalUsedStorage)

		estimatedfreelogicalstorage := humanize.IBytes(info.Statistics.EstimatedFreeLogicalStorage)
		estimatedlogicalstorage := humanize.IBytes(info.Statistics.EstimatedLogicalStorage)
		logicalstorage := humanize.IBytes(info.Statistics.LogicalStorage)
		logicalusedstorage := humanize.IBytes(info.Statistics.LogicalUsedStorage)
		installedphysicalstorage := humanize.IBytes(info.Statistics.InstalledPhysicalStorage)
		managedphysicalstorage := humanize.IBytes(info.Statistics.ManagedPhysicalStorage)
		// physicalusedstorageincludingparity :=  humanize.IBytes(info.Statistics.PhysicalUsedStorageIncludingParity)

		version := "n/a"
		if info.MinVersionInCluster != "" {
			version = info.MinVersionInCluster
		}
		short := []string{
			partition,
			version,
			health,
			fmt.Sprintf("%d", numinactivenodes),
			fmt.Sprintf("%d/%d/%d", numdegradedvolumes, numnotavailablevolumes, numreadonlyvolumes),
			installedphysicalstorage + "/" + managedphysicalstorage,
			effectivephysicalstorage + "/" + freephysicalstorage + "/" + physicalusedstorage,
			logicalstorage + "/" + logicalusedstorage,
			estimatedlogicalstorage + "/" + estimatedfreelogicalstorage,
			compressionratio,
		}
		rows = append(rows, short)
	}

	return header, rows, nil
}
