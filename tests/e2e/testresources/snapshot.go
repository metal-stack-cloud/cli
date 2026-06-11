package testresources

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Snapshot1 = func() *apiv1.Snapshot {
		return &apiv1.Snapshot{
			Uuid:             "8f4c20b8-406a-4952-b91c-73c38dbed112",
			Name:             "prod-db-backup-daily",
			Description:      "Daily automated backup of the production database cluster",
			Project:          "project-99",
			Partition:        "us-east-1a",
			StorageClass:     "premium-nvme",
			Size:             53687091200,
			Usage:            42949672960,
			State:            "READY",
			SourceVolumeUuid: "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
			SourceVolumeName: "prod-db-volume-01",
			ReplicaCount:     3,
			PrimaryNodeUuid:  "node-storage-prd-04",
			Retention:        durationpb.New(7 * 24 * time.Hour),
			CreatedAt:        timestamppb.New(e2e.TimeBubbleStartTime().Add(-2 * time.Hour)),
			Statistics: &apiv1.SnapshotStatistics{
				PhysicalCapacity:      1000000000,
				PhysicalOwnedCapacity: 250000000,
				PhysicalOwnedMemory:   16777216,
				PhysicalMemory:        67108864,
				UserWritten:           500000000,
			},
		}
	}
	Snapshot2 = func() *apiv1.Snapshot {
		return &apiv1.Snapshot{
			Uuid:             "3f4g20b8-406a-4952-b91c-73c38dbed113",
			Name:             "prod-db-backup-hour",
			Description:      "Hourly automated backup of the production database cluster",
			Project:          "project-99",
			Partition:        "us-east-1a",
			StorageClass:     "premium-nvme",
			Size:             53687091200,
			Usage:            42949672960,
			State:            "READY",
			SourceVolumeUuid: "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
			SourceVolumeName: "prod-db-volume-01",
			ReplicaCount:     3,
			PrimaryNodeUuid:  "node-storage-prd-04",
			Retention:        durationpb.New(7 * 24 * time.Hour),
			CreatedAt:        timestamppb.New(e2e.TimeBubbleStartTime().Add(-2 * time.Hour)),
			Statistics: &apiv1.SnapshotStatistics{
				PhysicalCapacity:      1000000000,
				PhysicalOwnedCapacity: 250000000,
				PhysicalOwnedMemory:   16777216,
				PhysicalMemory:        67108864,
				UserWritten:           500000000,
			},
		}
	}
)
