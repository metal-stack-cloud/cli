package testresources

import (
	"time"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
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
	Machine1 = func() *adminv1.Machine {
		return &adminv1.Machine{
			Uuid:        "123e4567-e89b-12d3-a456-426614174000",
			Name:        "test-compute-node-01",
			Description: "A dummy machine instance for unit and integration testing.",
			Project:     "project-99",
			Image:       "ubuntu-22.04-lts",
			Size:        "c1-large-x86",
			Hostname:    "node-01.internal.net",
			UserData:    "data",
			Role:        "machine",
			Creator:     "admin-user@company.com",
			Created:     timestamppb.New(e2e.TimeBubbleStartTime().Add(-24 * time.Hour)),
			Partition:   "us-east-1a",
			Rack:        "rack-b4-u12",
			State:       "RUNNING",
			Liveliness:  "ALIVE",
			Tags: []string{
				"environment:testing",
				"tier:frontend",
				"ephemeral",
			},
			MachineNetworks: []*adminv1.MachineNetwork{
				{
					Network:             "b8a3f290-7c1a-4d3b-8f1a-9c4b5e6f7a8b",
					Prefixes:            []string{"10.0.0.0/24", "2001:db8:1::/64"},
					Ips:                 []string{"10.0.0.15", "2001:db8:1::15"},
					DestinationPrefixes: []string{"0.0.0.0/0", "::/0"},
					NetworkType:         "overlay",
					Vrf:                 1001,
					Asn:                 65500,
				},
			},
			Vpn: &adminv1.VPN{
				Address: "10.0.0.5",
				Authkey: "tskey-auth-k123456CNTRL-al9dkf205jfh",
			},
		}
	}
	Machine2 = func() *adminv1.Machine {
		return &adminv1.Machine{
			Uuid:        "ea01dc1e-a349-4283-b32d-b432733a6d06",
			Name:        "test-compute-node-02",
			Description: "A dummy machine instance for unit and integration testing.",
			Project:     "project-01",
			Image:       "ubuntu-22.04-lts",
			Size:        "n1-large-x86",
			Hostname:    "node-01.internal.net",
			UserData:    "data",
			Role:        "machine",
			Creator:     "admin-user@company.com",
			Created:     timestamppb.New(e2e.TimeBubbleStartTime()),
			Partition:   "muc-1",
			Rack:        "rack-b4-u15",
			State:       "RUNNING",
			Liveliness:  "ALIVE",
			Tags: []string{
				"environment:prod",
			},
			MachineNetworks: []*adminv1.MachineNetwork{
				{
					Network:             "b8a3f290-7c1a-4d3b-8f1a-9c4b5e6f7a8b",
					Prefixes:            []string{"10.0.0.0/24", "2001:db8:1::/64"},
					Ips:                 []string{"10.0.0.15", "2001:db8:1::15"},
					DestinationPrefixes: []string{"0.0.0.0/0", "::/0"},
					NetworkType:         "overlay",
					Vrf:                 1001,
					Asn:                 65500,
				},
			},
			Vpn: &adminv1.VPN{
				Address: "10.0.0.5",
				Authkey: "hskey-auth-k123456CNTRL-al9dkf2g5jfh",
			},
		}
	}
)
