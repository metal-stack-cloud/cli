package cmd

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp/cmpopts"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	v1 "github.com/metal-stack-cloud/cli/cmd/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	cluster1 = func() *apiv1.Cluster {
		return &apiv1.Cluster{
			Uuid:       "6c631ff1-9038-4ad0-b75e-3ea173b7cdb1",
			Name:       "cluster1",
			Project:    "a",
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
				KubernetesAutoupdate:   pointer.Pointer(true),
				MachineimageAutoupdate: pointer.Pointer(false),
				TimeWindow: &apiv1.MaintenanceTimeWindow{
					Begin: &apiv1.Time{
						Hour:     18,
						Minute:   30,
						Timezone: "UTC",
					},
					Duration: durationpb.New(1 * time.Hour),
				},
			},
			Tenant:    "metal-stack",
			CreatedAt: timestamppb.New(testTime),
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
				LastErrors:            nil,
			},
			Purpose:    pointer.Pointer("evaluation"),
			Monitoring: &apiv1.ClusterMonitoring{},
		}
	}
	cluster2 = func() *apiv1.Cluster {
		return &apiv1.Cluster{
			Uuid:       "0c538734-c469-46a0-8efd-98e439d4dc8a",
			Name:       "cluster2",
			Project:    "a",
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
			Tenant:      "metal-stack",
			CreatedAt:   timestamppb.New(testTime),
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
			Purpose:    pointer.Pointer("production"),
			Monitoring: &apiv1.ClusterMonitoring{},
		}
	}
)

func Test_ClusterCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.Cluster]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.Cluster) []string {
				return []string{"cluster", "list", "--project", "a"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("List", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceListRequest{
							Project: "a",
						})).Return(&connect.Response[apiv1.ClusterServiceListResponse]{
							Msg: &apiv1.ClusterServiceListResponse{
								Clusters: []*apiv1.Cluster{
									cluster2(),
									cluster1(),
								},
							},
						}, nil)
					},
				},
			},
			Want: []*apiv1.Cluster{
				cluster1(),
				cluster2(),
			},
			WantTable: pointer.Pointer(`
TENANT       PROJECT  ID                                    NAME      PARTITION    VERSION  SIZE   AGE  
72%   metal-stack  a        6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  partition-a  1.25.10  1 - 3  now  
100%  metal-stack  a        0c538734-c469-46a0-8efd-98e439d4dc8a  cluster2  partition-b  1.27.9   3 - 6  now
`),
			WantWideTable: pointer.Pointer(`
ID                                    TENANT       PROJECT  NAME      PARTITION    PURPOSE     VERSION  OPERATION   PROGRESS          API  CONTROL  NODES  SYS  SIZE   AGE  
6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  metal-stack  a        cluster1  partition-a  evaluation  1.25.10  Processing  72% [Reconcile]   ✔    ✔        ✗      ✔    1 - 3  now  
0c538734-c469-46a0-8efd-98e439d4dc8a  metal-stack  a        cluster2  partition-b  production  1.27.9   Succeeded   100% [Reconcile]  ✔    ✔        ✔      ✔    3 - 6  now
`),
			Template: pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 a
0c538734-c469-46a0-8efd-98e439d4dc8a a
			`),
			WantMarkdown: pointer.Pointer(`
|      | TENANT      | PROJECT | ID                                   | NAME     | PARTITION   | VERSION | SIZE  | AGE |
|------|-------------|---------|--------------------------------------|----------|-------------|---------|-------|-----|
| 72%  | metal-stack | a       | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | partition-a | 1.25.10 | 1 - 3 | now |
| 100% | metal-stack | a       | 0c538734-c469-46a0-8efd-98e439d4dc8a | cluster2 | partition-b | 1.27.9  | 3 - 6 | now |
`),
		},
		{
			Name: "apply",
			Cmd: func(want []*apiv1.Cluster) []string {
				return appendFromFileCommonArgs("cluster", "apply")
			},
			FsMocks: func(fs afero.Fs, want []*apiv1.Cluster) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Create", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.ClusterResponseToCreate(cluster1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
							Cluster: cluster1(),
						}), nil)
						m.On("Create", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.ClusterResponseToCreate(cluster2())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
							Cluster: cluster2(),
						}), nil)
						// FIXME: API does not return a conflict when already exists, so the update functionality does not work!
					},
				},
			},
			Want: []*apiv1.Cluster{
				cluster1(),
				cluster2(),
			},
		},
		{
			Name: "update from file",
			Cmd: func(want []*apiv1.Cluster) []string {
				return appendFromFileCommonArgs("cluster", "update")
			},
			FsMocks: func(fs afero.Fs, want []*apiv1.Cluster) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.ClusterResponseToUpdate(cluster1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceUpdateResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: []*apiv1.Cluster{
				cluster1(),
			},
		},
		{
			Name: "create from file",
			Cmd: func(want []*apiv1.Cluster) []string {
				return appendFromFileCommonArgs("cluster", "create")
			},
			FsMocks: func(fs afero.Fs, want []*apiv1.Cluster) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Create", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.ClusterResponseToCreate(cluster1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: []*apiv1.Cluster{
				cluster1(),
			},
		},
		{
			Name: "delete from file",
			Cmd: func(want []*apiv1.Cluster) []string {
				return appendFromFileCommonArgs("cluster", "delete")
			},
			FsMocks: func(fs afero.Fs, want []*apiv1.Cluster) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceDeleteRequest{
							Uuid:    cluster1().Uuid,
							Project: cluster1().Project,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceDeleteResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: []*apiv1.Cluster{
				cluster1(),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_SingleResult(t *testing.T) {
	tests := []*Test[*apiv1.Cluster]{
		{
			Name: "describe",
			Cmd: func(want *apiv1.Cluster) []string {
				return []string{"cluster", "describe", "--project", want.Project, want.Uuid}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
							Project: cluster1().Project,
							Uuid:    cluster1().Uuid,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceGetResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: cluster1(),
			WantTable: pointer.Pointer(`
TENANT       PROJECT  ID                                    NAME      PARTITION    VERSION  SIZE   AGE  
72%  metal-stack  a        6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  partition-a  1.25.10  1 - 3  now
`),
			WantWideTable: pointer.Pointer(`
ID                                    TENANT       PROJECT  NAME      PARTITION    PURPOSE     VERSION  OPERATION   PROGRESS         API  CONTROL  NODES  SYS  SIZE   AGE  
6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  metal-stack  a        cluster1  partition-a  evaluation  1.25.10  Processing  72% [Reconcile]  ✔    ✔        ✗      ✔    1 - 3  now
`),
			Template: pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 a
			`),
			WantMarkdown: pointer.Pointer(`
|     | TENANT      | PROJECT | ID                                   | NAME     | PARTITION   | VERSION | SIZE  | AGE |
|-----|-------------|---------|--------------------------------------|----------|-------------|---------|-------|-----|
| 72% | metal-stack | a       | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | partition-a | 1.25.10 | 1 - 3 | now |
`),
		},
		{
			Name: "delete",
			Cmd: func(want *apiv1.Cluster) []string {
				return []string{"cluster", "rm", "--project", want.Project, want.Uuid, "--skip-security-prompts"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceDeleteRequest{
							Project: cluster1().Project,
							Uuid:    cluster1().Uuid,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceDeleteResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: cluster1(),
		},
		{
			Name: "create",
			Cmd: func(want *apiv1.Cluster) []string {
				args := []string{"cluster", "create",
					"--project", want.Project,
					"--partition", want.Partition,
					"--name", want.Name,
					"--kubernetes-version", want.Kubernetes.Version,
					"--maintenance-duration", want.Maintenance.TimeWindow.Duration.AsDuration().String(),
					"--maintenance-hour", strconv.Itoa(int(want.Maintenance.TimeWindow.Begin.Hour)), // nolint:gosec
					"--maintenance-minute", strconv.Itoa(int(want.Maintenance.TimeWindow.Begin.Minute)), // nolint:gosec
					"--maintenance-timezone", want.Maintenance.TimeWindow.Begin.Timezone,
					"--worker-group", want.Workers[0].Name,
					"--worker-min", strconv.Itoa(int(want.Workers[0].Minsize)), // nolint:gosec
					"--worker-max", strconv.Itoa(int(want.Workers[0].Maxsize)), // nolint:gosec
					"--worker-max-surge", strconv.Itoa(int(want.Workers[0].Maxsurge)), // nolint:gosec
					"--worker-max-unavailable", strconv.Itoa(int(want.Workers[0].Maxunavailable)), // nolint:gosec
					"--worker-type", want.Workers[0].MachineType,
				}
				AssertExhaustiveArgs(t, args, commonExcludedFileArgs()...)
				return args
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						req := cluster1()
						req.Maintenance.KubernetesAutoupdate = nil
						req.Maintenance.MachineimageAutoupdate = nil
						m.On("Create", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.ClusterResponseToCreate(req)), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: cluster1(),
		},
		{
			Name: "update",
			Cmd: func(want *apiv1.Cluster) []string {
				args := []string{"cluster", "update", want.Uuid,
					"--project", want.Project,
					"--kubernetes-version", want.Kubernetes.Version,
					"--maintenance-duration", want.Maintenance.TimeWindow.Duration.AsDuration().String(),
					"--maintenance-hour", strconv.Itoa(int(want.Maintenance.TimeWindow.Begin.Hour)), // nolint:gosec
					"--maintenance-minute", strconv.Itoa(int(want.Maintenance.TimeWindow.Begin.Minute)), // nolint:gosec
					"--maintenance-timezone", want.Maintenance.TimeWindow.Begin.Timezone,
					"--worker-group", want.Workers[0].Name,
					"--worker-min", strconv.Itoa(int(want.Workers[0].Minsize)), // nolint:gosec
					"--worker-max", strconv.Itoa(int(want.Workers[0].Maxsize)), // nolint:gosec
					"--worker-max-surge", strconv.Itoa(int(want.Workers[0].Maxsurge)), // nolint:gosec
					"--worker-max-unavailable", strconv.Itoa(int(want.Workers[0].Maxunavailable)), // nolint:gosec
					"--worker-type", want.Workers[0].MachineType,
				}
				exclude := append(commonExcludedFileArgs(), "remove-worker-group")
				AssertExhaustiveArgs(t, args, exclude...)
				return args
			},
			MockStdin: bytes.NewBufferString("y"),
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Cluster: func(m *mock.Mock) {
						m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
							Uuid:    cluster1().Uuid,
							Project: cluster1().Project,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceGetResponse{
							Cluster: cluster1(),
						}), nil)

						m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.ClusterResponseToUpdate(cluster1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.ClusterServiceUpdateResponse{
							Cluster: cluster1(),
						}), nil)
					},
				},
			},
			Want: cluster1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
