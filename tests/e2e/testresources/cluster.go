package testresources

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Cluster1 = func() *apiv1.Cluster {
		return &apiv1.Cluster{
			Uuid:       "6c631ff1-9038-4ad0-b75e-3ea173b7cdb1",
			Name:       "cluster1",
			Project:    Project1().Uuid,
			Partition:  "partition-a",
			Kubernetes: &apiv1.KubernetesSpec{Version: "1.25.10"},
			Workers: []*apiv1.Worker{
				{
					Name:           "group-0",
					MachineType:    "c1-xlarge-x86",
					Minsize:        1,
					Maxsize:        3,
					Maxsurge:       1,
					Maxunavailable: 0,
				},
			},
			Maintenance: &apiv1.Maintenance{
				KubernetesAutoupdate:   new(true),
				MachineimageAutoupdate: new(false),
				TimeWindow: &apiv1.MaintenanceTimeWindow{
					Begin: &apiv1.Time{
						Hour:     18,
						Minute:   30,
						Timezone: "UTC",
					},
					Duration: durationpb.New(1 * time.Hour),
				},
			},
			Tenant:    Tenant1().Login,
			CreatedAt: timestamppb.New(e2e.TimeBubbleStartTime()),
			UpdatedAt: &timestamppb.Timestamp{},
			DeletedAt: &timestamppb.Timestamp{},
			Status: &apiv1.ClusterStatus{
				Uuid:                  "6c631ff1-9038-4ad0-b75e-3ea173b7cdb1",
				Progress:              72,
				State:                 "Processing",
				Type:                  "Reconcile",
				ApiServerReady:        "True",
				ControlPlaneReady:     "True",
				NodesReady:            "False",
				SystemComponentsReady: "True",
				LastErrors: []*apiv1.ClusterStatusLastError{
					&apiv1.ClusterStatusLastError{
						Description:    "failed",
						TaskId:         new("someid"),
						LastUpdateTime: timestamppb.New(e2e.TimeBubbleStartTime()),
					},
				},
				Conditions: []*apiv1.ClusterStatusCondition{
					&apiv1.ClusterStatusCondition{
						Type:               "Ready",
						Status:             "True",
						Reason:             "AllNodesHealthy",
						StatusMessage:      "All cluster nodes are reporting healthy status",
						LastTransitionTime: timestamppb.New(e2e.TimeBubbleStartTime().Add(-2 * time.Minute)),
						LastUpdateTime:     timestamppb.New(e2e.TimeBubbleStartTime()),
					},
				},
			},
			Purpose: new("evaluation"),
			Monitoring: &apiv1.ClusterMonitoring{
				Username: "username",
				Password: "password",
				Endpoint: "endpoint",
			},
		}
	}
	Cluster2 = func() *apiv1.Cluster {
		return &apiv1.Cluster{
			Uuid:       "0c538734-c469-46a0-8efd-98e439d4dc8a",
			Name:       "cluster2",
			Project:    Project2().Uuid,
			Partition:  "partition-b",
			Kubernetes: &apiv1.KubernetesSpec{Version: "1.27.9"},
			Workers: []*apiv1.Worker{
				{
					Minsize: 1,
					Maxsize: 3,
				},
				{
					Minsize: 2,
					Maxsize: 3,
				},
			},
			Maintenance: &apiv1.Maintenance{},
			Tenant:      Tenant2().Login,
			CreatedAt:   timestamppb.New(e2e.TimeBubbleStartTime()),
			UpdatedAt:   &timestamppb.Timestamp{},
			DeletedAt:   &timestamppb.Timestamp{},
			Status: &apiv1.ClusterStatus{
				Uuid:                  "0c538734-c469-46a0-8efd-98e439d4dc8a",
				Progress:              100,
				State:                 "Succeeded",
				Type:                  "Reconcile",
				ApiServerReady:        "True",
				ControlPlaneReady:     "True",
				NodesReady:            "True",
				SystemComponentsReady: "True",
				LastErrors:            nil,
			},
			Purpose:    new("production"),
			Monitoring: &apiv1.ClusterMonitoring{},
		}
	}
)
