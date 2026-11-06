package testresources

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

var (
	Volume1 = func() *apiv1.Volume {
		return &apiv1.Volume{
			Uuid:               "bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
			Name:               "volume1",
			Project:            Project1().Uuid,
			Partition:          "partition-a",
			StorageClass:       "storageclass-a",
			Size:               1024,
			Usage:              42,
			State:              "Bound",
			AttachedTo:         nil,
			SourceSnapshotUuid: "",
			SourceSnapshotName: "",
			VolumeHandle:       "",
			NodeIps:            nil,
			RebuildProgress:    "",
			PrimaryNodeUuid:    "",
			QosPolicyUuid:      "",
			QosPolicyName:      "",
			ReplicaCount:       0,
			ProtectionState:    "",
			LogicalUsedStorage: 0,
			Statistics:         nil,
			ClusterName:        Cluster1().Name,
			ClusterId:          Cluster1().Uuid,
			Labels: []*apiv1.VolumeLabel{
				{
					Key:   "foo",
					Value: "bar",
				},
			},
		}
	}
	Volume2 = func() *apiv1.Volume {
		return &apiv1.Volume{
			Uuid:               "0372d029-1077-4e9b-b303-7d64ad5496fd",
			Name:               "volume2",
			Project:            Project2().Uuid,
			Partition:          "partition-a",
			StorageClass:       "storageclass-a",
			Size:               1024,
			Usage:              42,
			State:              "Bound",
			AttachedTo:         nil,
			SourceSnapshotUuid: "",
			SourceSnapshotName: "",
			VolumeHandle:       "",
			NodeIps:            nil,
			RebuildProgress:    "",
			PrimaryNodeUuid:    "",
			QosPolicyUuid:      "",
			QosPolicyName:      "",
			ReplicaCount:       0,
			ProtectionState:    "",
			LogicalUsedStorage: 0,
			Statistics:         nil,
			ClusterName:        Cluster2().Name,
			ClusterId:          Cluster2().Uuid,
			Labels: []*apiv1.VolumeLabel{
				{
					Key:   "bar",
					Value: "baz",
				},
			},
		}
	}
)
