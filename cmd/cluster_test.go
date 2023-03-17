package cmd

import (
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/google/go-cmp/cmp/cmpopts"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ClusterCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.Cluster]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.Cluster) []string {
				return []string{"cluster", "list", "--project", "a"}
			},
			APIMocks: &apitests.Apiv1MockFns{
				Cluster: func(m *mock.Mock) {
					m.On("List", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceListRequest{
						Project: "a",
					})).Return(&connect.Response[apiv1.ClusterServiceListResponse]{
						Msg: &apiv1.ClusterServiceListResponse{
							Clusters: []*apiv1.Cluster{
								{
									Uuid:    "1-2-3-4",
									Name:    "a",
									Project: "a",
									Kubernetes: &apiv1.KubernetesSpec{
										Version: "1.23.1",
									},
									Workers: []*apiv1.Worker{
										{
											Name:        "a",
											MachineType: "c1-large-x86",
											Minsize:     4,
											Maxsize:     7,
										},
									},
									Status: &apiv1.ClusterStatus{
										State: "Processing",
									},
									CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
								},
								{
									Uuid:    "5-6-7-8",
									Name:    "b",
									Project: "a",
									Kubernetes: &apiv1.KubernetesSpec{
										Version: "1.23.1",
									},
									Workers: []*apiv1.Worker{
										{
											Name:        "b",
											MachineType: "c1-large-x86",
											Minsize:     1,
											Maxsize:     3,
										},
									},
									Status: &apiv1.ClusterStatus{
										State: "Processing",
									},
									CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Minute)),
								},
							},
						},
					}, nil)
				},
			},
			Want: []*apiv1.Cluster{
				{
					Uuid:    "5-6-7-8",
					Name:    "b",
					Project: "a",
					Kubernetes: &apiv1.KubernetesSpec{
						Version: "1.23.1",
					},
					Workers: []*apiv1.Worker{
						{
							Name:        "b",
							MachineType: "c1-large-x86",
							Minsize:     1,
							Maxsize:     3,
						},
					},
					Status: &apiv1.ClusterStatus{
						State: "Processing",
					},
					CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Minute)),
				},
				{
					Uuid:    "1-2-3-4",
					Name:    "a",
					Project: "a",
					Kubernetes: &apiv1.KubernetesSpec{
						Version: "1.23.1",
					},
					Workers: []*apiv1.Worker{
						{
							Name:        "a",
							MachineType: "c1-large-x86",
							Minsize:     4,
							Maxsize:     7,
						},
					},
					Status: &apiv1.ClusterStatus{
						State: "Processing",
					},
					CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
				},
			},
			WantTable: pointer.Pointer(`
CLUSTERSTATUS   ID        NAME   PROJECT   KUBERNETES VERSION   NODES   UPTIME
●               5-6-7-8   b      a         1.23.1               1 - 3   1 minute ago
●               1-2-3-4   a      a         1.23.1               4 - 7   1 hour ago
`),
			WantWideTable: pointer.Pointer(`
CLUSTERSTATUS   ID        NAME   PROJECT   KUBERNETES VERSION   NODES   UPTIME
●               5-6-7-8   b      a         1.23.1               1 - 3   1 minute ago
●               1-2-3-4   a      a         1.23.1               4 - 7   1 hour ago
`),
			// 			Template: pointer.Pointer("{{ .status.state }} {{ .uuid }} {{ .name }} {{ .project }} {{ .kubernetes.version }} {{ .workers[0].minsize - .workers[0].maxsize}}"), // TODO How to do min and max size?
			// 			WantTemplate: pointer.Pointer(`
			// Processing 1-2-3-4 a a 1.23.1 4 - 7
			// Processing 5-6-7-8 b a 1.23.1 1 - 3
			// `),
			WantMarkdown: pointer.Pointer(`
| CLUSTERSTATUS |   ID    | NAME | PROJECT | KUBERNETES VERSION | NODES |    UPTIME    |
|---------------|---------|------|---------|--------------------|-------|--------------|
| ●             | 5-6-7-8 | b    | a       | 1.23.1             | 1 - 3 | 1 minute ago |
| ●             | 1-2-3-4 | a    | a       | 1.23.1             | 4 - 7 | 1 hour ago   |
`),
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
				return []string{"cluster", "describe", "--project", "proj-a", "1-2-3-4"}
			},
			APIMocks: &apitests.Apiv1MockFns{
				Cluster: func(m *mock.Mock) {
					m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
						Uuid:    "1-2-3-4",
						Project: "proj-a",
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.ClusterServiceGetResponse]{
						Msg: &apiv1.ClusterServiceGetResponse{
							Cluster: &apiv1.Cluster{
								Uuid:      "1-2-3-4",
								Name:      "name-a",
								Project:   "proj-a",
								Partition: "part-1",
								Kubernetes: &apiv1.KubernetesSpec{
									Version: "1.23.1",
								},
								Workers: []*apiv1.Worker{
									{
										Name:        "wogr-1",
										MachineType: "m1-mega",
										Minsize:     1,
										Maxsize:     3,
									},
								},
								Maintenance: &apiv1.Maintenance{
									TimeWindow: &apiv1.MaintenanceTimeWindow{
										Begin:    timestamppb.New(testTime),
										Duration: durationpb.New(1 * time.Hour),
									},
								},
								Status: &apiv1.ClusterStatus{
									State: "Processing",
								},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.Cluster{
				Uuid:      "1-2-3-4",
				Name:      "name-a",
				Project:   "proj-a",
				Partition: "part-1",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.23.1",
				},
				Workers: []*apiv1.Worker{
					{
						Name:        "wogr-1",
						MachineType: "m1-mega",
						Minsize:     1,
						Maxsize:     3,
					},
				},
				Maintenance: &apiv1.Maintenance{
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin:    timestamppb.New(testTime),
						Duration: durationpb.New(1 * time.Hour),
					},
				},
				Status: &apiv1.ClusterStatus{
					State: "Processing",
				},
			},
			// TODO What is the problem with the tables?
			// 			WantTable: pointer.Pointer(`
			// CLUSTERSTATUS   ID        NAME     PROJECT   KUBERNETES VERSION
			// ●               1-2-3-4   name-a   proj-a    1.23.1
			// `),
			// 			WantWideTable: pointer.Pointer(`
			// CLUSTERSTATUS   ID        NAME   PROJECT   KUBERNETESSPEC
			// Processing      1-2-3-4   a      a         1.23.1
			// `),
			// 			Template: pointer.Pointer("{{ .uuid }} {{ .name }} {{ .project }}"), //TODO How to work with status and kubernetes version here?
			// 			WantTemplate: pointer.Pointer(`
			// 1-2-3-4 a a
			// `),
			// 			WantMarkdown: pointer.Pointer(`
			// | CLUSTERSTATUS |   ID    | NAME | PROJECT | KUBERNETESSPEC |
			// |---------------|---------|------|---------|----------------|
			// | Processing    | 1-2-3-4 | a    | a       | 1.23.1         |
			// `),
		},
		{
			Name: "delete",
			Cmd: func(want *apiv1.Cluster) []string {
				return []string{"cluster", "delete", "--project", "a", "1-2-3-4"}
			},
			APIMocks: &apitests.Apiv1MockFns{
				Cluster: func(m *mock.Mock) {
					m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceDeleteRequest{
						Project: "a",
						Uuid:    "1-2-3-4",
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.ClusterServiceDeleteResponse]{
						Msg: &apiv1.ClusterServiceDeleteResponse{
							Cluster: &apiv1.Cluster{
								Uuid:    "1-2-3-4",
								Project: "a",
								Name:    "a",
								Kubernetes: &apiv1.KubernetesSpec{
									Version: "1.23.1",
								},
								Workers: []*apiv1.Worker{
									{
										MachineType: "c1-large",
										Minsize:     1,
										Maxsize:     3,
									},
								},
								Status: &apiv1.ClusterStatus{
									State: "Processing",
								},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.Cluster{
				Uuid:    "1-2-3-4",
				Project: "a",
				Name:    "a",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.23.1",
				},
				Workers: []*apiv1.Worker{
					{
						MachineType: "c1-large",
						Minsize:     1,
						Maxsize:     3,
					},
				},
				Status: &apiv1.ClusterStatus{
					State: "Processing",
				},
			},
		},
		{
			Name: "create",
			Cmd: func(want *apiv1.Cluster) []string {
				args := []string{"cluster", "create", "--name", "abc", "--project", "a", "--partition", "part-1", "--kubernetes", "1.1.1", "--workername", "a-worker", "--machinetype", "ma-large", "--minsize", "1", "--maxsize", "3", "--maintenancebegin", "04:30pm", "--maintenanceduration", "1h"}
				AssertExhaustiveArgs(t, args, "file")
				return args
			},
			APIMocks: &apitests.Apiv1MockFns{
				Cluster: func(m *mock.Mock) {
					m.On("Create", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
						Project:   "a",
						Name:      "abc",
						Partition: "part-1",
						Kubernetes: &apiv1.KubernetesSpec{
							Version: "1.1.1",
						},
						Workers: []*apiv1.Worker{
							{
								Name:        "a-worker",
								MachineType: "ma-large",
								Minsize:     1,
								Maxsize:     3,
							},
						},
						Maintenance: &apiv1.Maintenance{
							TimeWindow: &apiv1.MaintenanceTimeWindow{
								Begin: &timestamppb.Timestamp{
									Seconds: 59400,
								},
								Duration: &durationpb.Duration{
									Seconds: 3600,
								},
							},
						},
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.ClusterServiceCreateResponse]{
						Msg: &apiv1.ClusterServiceCreateResponse{
							Cluster: &apiv1.Cluster{
								Project:   "a",
								Name:      "abc",
								Partition: "part-1",
								Kubernetes: &apiv1.KubernetesSpec{
									Version: "1.1.1",
								},
								Workers: []*apiv1.Worker{
									{
										Name:        "a-worker",
										MachineType: "ma-large",
										Minsize:     1,
										Maxsize:     3,
									},
								},
								Maintenance: &apiv1.Maintenance{
									TimeWindow: &apiv1.MaintenanceTimeWindow{
										Begin: &timestamppb.Timestamp{
											Seconds: 59400,
										},
										Duration: &durationpb.Duration{
											Seconds: 3600,
										},
									},
								},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.Cluster{
				Project:   "a",
				Name:      "abc",
				Partition: "part-1",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.1.1",
				},
				Workers: []*apiv1.Worker{
					{
						Name:        "a-worker",
						MachineType: "ma-large",
						Minsize:     1,
						Maxsize:     3,
					},
				},
				Maintenance: &apiv1.Maintenance{
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin: &timestamppb.Timestamp{
							Seconds: 59400,
						},
						Duration: &durationpb.Duration{
							Seconds: 3600,
						},
					},
				},
			},
		},
		{
			Name: "create from file",
			Cmd: func(want *apiv1.Cluster) []string {
				return []string{"cluster", "create", "-f", "/file.yaml"}
			},
			FsMocks: func(fs afero.Fs, want *apiv1.Cluster) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshal(t, want), 0755))
			},
			APIMocks: &apitests.Apiv1MockFns{
				Cluster: func(m *mock.Mock) {
					m.On("Create", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
						Project:   "a",
						Name:      "abc",
						Partition: "part-1",
						Kubernetes: &apiv1.KubernetesSpec{
							Version: "1.1.1",
						},
						Workers: []*apiv1.Worker{
							{
								Name:        "a-worker",
								MachineType: "large",
								Minsize:     1,
								Maxsize:     3,
							},
						},
						Maintenance: &apiv1.Maintenance{
							TimeWindow: &apiv1.MaintenanceTimeWindow{
								Begin:    timestamppb.New(testTime),
								Duration: durationpb.New(1 * time.Hour),
							},
						},
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.ClusterServiceCreateResponse]{
						Msg: &apiv1.ClusterServiceCreateResponse{
							Cluster: &apiv1.Cluster{
								Project:   "a",
								Name:      "abc",
								Partition: "part-1",
								Kubernetes: &apiv1.KubernetesSpec{
									Version: "1.1.1",
								},
								Workers: []*apiv1.Worker{
									{
										Name:        "a-worker",
										MachineType: "large",
										Minsize:     1,
										Maxsize:     3,
									},
								},
								Maintenance: &apiv1.Maintenance{
									TimeWindow: &apiv1.MaintenanceTimeWindow{
										Begin:    timestamppb.New(testTime),
										Duration: durationpb.New(1 * time.Hour),
									},
								},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.Cluster{
				Project:   "a",
				Name:      "abc",
				Partition: "part-1",
				Kubernetes: &apiv1.KubernetesSpec{
					Version: "1.1.1",
				},
				Workers: []*apiv1.Worker{
					{
						Name:        "a-worker",
						MachineType: "large",
						Minsize:     1,
						Maxsize:     3,
					},
				},
				Maintenance: &apiv1.Maintenance{
					TimeWindow: &apiv1.MaintenanceTimeWindow{
						Begin:    timestamppb.New(testTime),
						Duration: durationpb.New(1 * time.Hour),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

// package cmd
