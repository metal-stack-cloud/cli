package api_e2e

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_ClusterCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceListResponse, []*apiv1.Cluster]{
		{
			Name:    "list",
			CmdArgs: []string{"cluster", "list", "--project", testresources.Cluster2().Project},

			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceListRequest{
								Project: testresources.Cluster2().Project,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceListResponse{
								Clusters: []*apiv1.Cluster{
									testresources.Cluster2(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT         PARTITION    VERSION  SIZE   AGE  
            100%      0c538734-c469-46a0-8efd-98e439d4dc8a  cluster2  c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent.com  partition-b  1.27.9   3 - 6  now
`),
			WantWideTable: new(`
            PROGRESS          ID                                    NAME      PROJECT                               TENANT         PARTITION    VERSION  SIZE   AGE  PURPOSE     OPERATION  API  CONTROL  NODES  SYS  
            100% [Reconcile]  0c538734-c469-46a0-8efd-98e439d4dc8a  cluster2  c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent.com  partition-b  1.27.9   3 - 6  now  production  Succeeded  ✔    ✔        ✔      ✔
`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
0c538734-c469-46a0-8efd-98e439d4dc8a c40ad996-e1fd-4511-a7bf-418219cb8d67
			`),
			WantMarkdown: new(`
            | PROGRESS | ID                                   | NAME     | PROJECT                              | TENANT        | PARTITION   | VERSION | SIZE  | AGE |
            |----------|--------------------------------------|----------|--------------------------------------|---------------|-------------|---------|-------|-----|
            | 100%     | 0c538734-c469-46a0-8efd-98e439d4dc8a | cluster2 | c40ad996-e1fd-4511-a7bf-418219cb8d67 | x-cellent.com | partition-b | 1.27.9  | 3 - 6 | now |
`)},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceListResponse, []*apiv1.Cluster]{
		{
			Name:    "apply from file",
			CmdArgs: append([]string{"cluster", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Cluster1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Cluster: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
									Name:        testresources.Cluster1().Name,
									Project:     testresources.Cluster1().Project,
									Partition:   testresources.Cluster1().Partition,
									Kubernetes:  testresources.Cluster1().Kubernetes,
									Workers:     testresources.Cluster1().Workers,
									Maintenance: testresources.Cluster1().Maintenance,
								})).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
									Cluster: testresources.Cluster1(),
								}), nil)
							},
						},
					}}),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
			`),
			WantMarkdown: new(`
            | PROGRESS | ID                                   | NAME     | PROJECT                              | TENANT      | PARTITION   | VERSION | SIZE  | AGE |
            |----------|--------------------------------------|----------|--------------------------------------|-------------|-------------|---------|-------|-----|
            | 72%      | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | partition-a | 1.25.10 | 1 - 3 | now |
			`),
		},
		{
			Name:    "apply many from file",
			CmdArgs: append([]string{"cluster", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Cluster1(), testresources.Cluster2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Cluster: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
									Name:        testresources.Cluster2().Name,
									Project:     testresources.Cluster2().Project,
									Partition:   testresources.Cluster2().Partition,
									Kubernetes:  testresources.Cluster2().Kubernetes,
									Workers:     testresources.Cluster2().Workers,
									Maintenance: testresources.Cluster2().Maintenance,
								})).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
									Cluster: testresources.Cluster2(),
								}), nil)
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
									Name:        testresources.Cluster1().Name,
									Project:     testresources.Cluster1().Project,
									Partition:   testresources.Cluster1().Partition,
									Kubernetes:  testresources.Cluster1().Kubernetes,
									Workers:     testresources.Cluster1().Workers,
									Maintenance: testresources.Cluster1().Maintenance,
								})).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
									Cluster: testresources.Cluster1(),
								}), nil)
								// FIXME: API does not return a conflict when already exists, so the update functionality does not work!
							},
						},
					}}),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT         PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack    partition-a  1.25.10  1 - 3  now  
            100%      0c538734-c469-46a0-8efd-98e439d4dc8a  cluster2  c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent.com  partition-b  1.27.9   3 - 6  now
			`),
			WantMarkdown: new(`
            | PROGRESS | ID                                   | NAME     | PROJECT                              | TENANT        | PARTITION   | VERSION | SIZE  | AGE |
            |----------|--------------------------------------|----------|--------------------------------------|---------------|-------------|---------|-------|-----|
            | 72%      | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack   | partition-a | 1.25.10 | 1 - 3 | now |
            | 100%     | 0c538734-c469-46a0-8efd-98e439d4dc8a | cluster2 | c40ad996-e1fd-4511-a7bf-418219cb8d67 | x-cellent.com | partition-b | 1.27.9  | 3 - 6 | now |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceListResponse, *apiv1.Cluster]{
		{
			Name: "update",
			CmdArgs: []string{"cluster", "update", testresources.Cluster1().Uuid,
				"--project", testresources.Cluster1().Project,
				"--kubernetes-version", testresources.Cluster1().Kubernetes.Version,
				"--maintenance-duration", testresources.Cluster1().Maintenance.TimeWindow.Duration.AsDuration().String(),
				"--maintenance-hour", strconv.Itoa(int(testresources.Cluster1().Maintenance.TimeWindow.Begin.Hour)), // nolint:gosec
				"--maintenance-minute", strconv.Itoa(int(testresources.Cluster1().Maintenance.TimeWindow.Begin.Minute)), // nolint:gosec
				"--maintenance-timezone", testresources.Cluster1().Maintenance.TimeWindow.Begin.Timezone,
				"--worker-group", testresources.Cluster1().Workers[0].Name,
				"--worker-min", strconv.Itoa(int(testresources.Cluster1().Workers[0].Minsize)), // nolint:gosec
				"--worker-max", strconv.Itoa(int(testresources.Cluster1().Workers[0].Maxsize)), // nolint:gosec
				"--worker-max-surge", strconv.Itoa(int(testresources.Cluster1().Workers[0].Maxsurge)), // nolint:gosec
				"--worker-max-unavailable", strconv.Itoa(int(testresources.Cluster1().Workers[0].Maxunavailable)), // nolint:gosec
				"--worker-type", testresources.Cluster1().Workers[0].MachineType,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				MockStdin: bytes.NewBufferString("y"),
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
								Uuid:    testresources.Cluster1().Uuid,
								Project: testresources.Cluster1().Project,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceGetResponse{
								Cluster: testresources.Cluster1(),
							}), nil)

							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceUpdateRequest{
								Uuid:        testresources.Cluster1().Uuid,
								Project:     testresources.Cluster1().Project,
								Kubernetes:  testresources.Cluster1().Kubernetes,
								Maintenance: testresources.Cluster1().Maintenance,
								Workers: []*apiv1.WorkerUpdate{
									&apiv1.WorkerUpdate{
										Name:           testresources.Cluster1().Workers[0].Name,
										MachineType:    &testresources.Cluster1().Workers[0].MachineType,
										Minsize:        &testresources.Cluster1().Workers[0].Minsize,
										Maxsize:        &testresources.Cluster1().Workers[0].Maxsize,
										Maxsurge:       &testresources.Cluster1().Workers[0].Maxsurge,
										Maxunavailable: &testresources.Cluster1().Workers[0].Maxunavailable,
									},
								},
							})).Return(connect.NewResponse(&apiv1.ClusterServiceUpdateResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Cluster1(),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
			`),
			WantWideTable: new(`
            PROGRESS         ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  PURPOSE     OPERATION   API  CONTROL  NODES  SYS  
            72% [Reconcile]  6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now  evaluation  Processing  ✔    ✔        ✗      ✔
			`),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"cluster", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Cluster1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Cluster: func(m *mock.Mock) {
								m.On("Update", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceUpdateRequest{
									Uuid:        testresources.Cluster1().Uuid,
									Project:     testresources.Cluster1().Project,
									Kubernetes:  testresources.Cluster1().Kubernetes,
									Maintenance: testresources.Cluster1().Maintenance,
									Workers: []*apiv1.WorkerUpdate{
										&apiv1.WorkerUpdate{
											Name:           testresources.Cluster1().Workers[0].Name,
											MachineType:    &testresources.Cluster1().Workers[0].MachineType,
											Minsize:        &testresources.Cluster1().Workers[0].Minsize,
											Maxsize:        &testresources.Cluster1().Workers[0].Maxsize,
											Maxsurge:       &testresources.Cluster1().Workers[0].Maxsurge,
											Maxunavailable: &testresources.Cluster1().Workers[0].Maxunavailable,
										},
									},
								})).Return(connect.NewResponse(&apiv1.ClusterServiceUpdateResponse{
									Cluster: testresources.Cluster1(),
								}), nil)
							},
						},
					},
				}),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
					`),
			WantWideTable: new(`
            PROGRESS         ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  PURPOSE     OPERATION   API  CONTROL  NODES  SYS  
            72% [Reconcile]  6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now  evaluation  Processing  ✔    ✔        ✗      ✔
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceListResponse, *apiv1.Cluster]{
		{
			Name: "create",
			CmdArgs: []string{"cluster", "create",
				"--project", testresources.Cluster1().Project,
				"--partition", testresources.Cluster1().Partition,
				"--name", testresources.Cluster1().Name,
				"--kubernetes-version", testresources.Cluster1().Kubernetes.Version,
				"--maintenance-duration", testresources.Cluster1().Maintenance.TimeWindow.Duration.AsDuration().String(),
				"--maintenance-hour", strconv.Itoa(int(testresources.Cluster1().Maintenance.TimeWindow.Begin.Hour)), // nolint:gosec
				"--maintenance-minute", strconv.Itoa(int(testresources.Cluster1().Maintenance.TimeWindow.Begin.Minute)), // nolint:gosec
				"--maintenance-timezone", testresources.Cluster1().Maintenance.TimeWindow.Begin.Timezone,
				"--worker-group", testresources.Cluster1().Workers[0].Name,
				"--worker-min", strconv.Itoa(int(testresources.Cluster1().Workers[0].Minsize)), // nolint:gosec
				"--worker-max", strconv.Itoa(int(testresources.Cluster1().Workers[0].Maxsize)), // nolint:gosec
				"--worker-max-surge", strconv.Itoa(int(testresources.Cluster1().Workers[0].Maxsurge)), // nolint:gosec
				"--worker-max-unavailable", strconv.Itoa(int(testresources.Cluster1().Workers[0].Maxunavailable)), // nolint:gosec
				"--worker-type", testresources.Cluster1().Workers[0].MachineType,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
								Name:       testresources.Cluster1().Name,
								Project:    testresources.Cluster1().Project,
								Partition:  testresources.Cluster1().Partition,
								Kubernetes: testresources.Cluster1().Kubernetes,
								Workers:    testresources.Cluster1().Workers,
								Maintenance: &apiv1.Maintenance{
									TimeWindow: &apiv1.MaintenanceTimeWindow{
										Begin:    testresources.Cluster1().Maintenance.TimeWindow.Begin,
										Duration: testresources.Cluster1().Maintenance.TimeWindow.Duration,
									},
								},
							})).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
			PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
			`),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"cluster", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Cluster1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Cluster: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceCreateRequest{
									Name:        testresources.Cluster1().Name,
									Project:     testresources.Cluster1().Project,
									Partition:   testresources.Cluster1().Partition,
									Kubernetes:  testresources.Cluster1().Kubernetes,
									Workers:     testresources.Cluster1().Workers,
									Maintenance: testresources.Cluster1().Maintenance,
								})).Return(connect.NewResponse(&apiv1.ClusterServiceCreateResponse{
									Cluster: testresources.Cluster1(),
								}), nil)
							},
						},
					},
				}),

			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
			`),
			WantWideTable: new(`
            PROGRESS         ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  PURPOSE     OPERATION   API  CONTROL  NODES  SYS  
            72% [Reconcile]  6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now  evaluation  Processing  ✔    ✔        ✗      ✔
			`),
			WantMarkdown: new(`
            | PROGRESS | ID                                   | NAME     | PROJECT                              | TENANT      | PARTITION   | VERSION | SIZE  | AGE |
            |----------|--------------------------------------|----------|--------------------------------------|-------------|-------------|---------|-------|-----|
            | 72%      | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | partition-a | 1.25.10 | 1 - 3 | now |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceListResponse, *apiv1.Cluster]{
		{
			Name:    "delete",
			CmdArgs: []string{"cluster", "rm", "--project", testresources.Cluster1().Project, testresources.Cluster1().Uuid, "--skip-security-prompts"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceDeleteRequest{
								Project: testresources.Cluster1().Project,
								Uuid:    testresources.Cluster1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceDeleteResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Cluster1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"cluster", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Cluster1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Cluster: func(m *mock.Mock) {
								m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceDeleteRequest{
									Uuid:    testresources.Cluster1().Uuid,
									Project: testresources.Cluster1().Project,
								})).Return(connect.NewResponse(&apiv1.ClusterServiceDeleteResponse{
									Cluster: testresources.Cluster1(),
								}), nil)
							},
						},
					}}),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceListResponse, *apiv1.Cluster]{
		{
			Name:    "describe",
			CmdArgs: []string{"cluster", "describe", "--project", testresources.Cluster1().Project, testresources.Cluster1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
								Project: testresources.Cluster1().Project,
								Uuid:    testresources.Cluster1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceGetResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Cluster1(),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now
`),
			WantWideTable: new(`
            PROGRESS         ID                                    NAME      PROJECT                               TENANT       PARTITION    VERSION  SIZE   AGE  PURPOSE     OPERATION   API  CONTROL  NODES  SYS  
            72% [Reconcile]  6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  partition-a  1.25.10  1 - 3  now  evaluation  Processing  ✔    ✔        ✗      ✔
`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 c40ad996-e1fd-4511-a7bf-418219cb8d95
			`),
			WantMarkdown: new(`
            | PROGRESS | ID                                   | NAME     | PROJECT                              | TENANT      | PARTITION   | VERSION | SIZE  | AGE |
            |----------|--------------------------------------|----------|--------------------------------------|-------------|-------------|---------|-------|-----|
            | 72%      | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | partition-a | 1.25.10 | 1 - 3 | now |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_ExecConfig(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceGetCredentialsResponse, string]{
		{
			Name:    "get exec-config",
			CmdArgs: []string{"cluster", "exec-config", "--project", testresources.Cluster1().Project, "--expiration", "9h", testresources.Cluster1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("GetCredentials", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceGetCredentialsRequest{
								Project:    testresources.Cluster1().Project,
								Uuid:       testresources.Cluster1().Uuid,
								Expiration: durationpb.New(9 * time.Hour),
							})).Return(connect.NewResponse(&apiv1.ClusterServiceGetCredentialsResponse{
								Kubeconfig: `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://example.com
  name: test
contexts:
- context:
    cluster: test
    user: test
  name: test
current-context: test
users:
- name: test
  user:
    token: fake
`,
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(`
{
  "kind": "ExecCredential",
  "apiVersion": "client.authentication.k8s.io/v1",
  "spec": {
    "interactive": false
  },
  "status": {
    "expirationTimestamp": "2000-01-01T09:00:00Z"
  }
}
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_KubeConfig(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceGetCredentialsResponse, string]{
		{
			Name: "get kubeconfig",
			CmdArgs: []string{"cluster", "kubeconfig", testresources.Cluster1().Uuid,
				"--project", testresources.Cluster1().Project, "--expiration", "9h", "--print-only=true", "--auth-type=certs"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
								Project: testresources.Cluster1().Project,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceGetResponse{
								Project: testresources.Project1(),
							}), nil)
						},
						Cluster: func(m *mock.Mock) {
							m.On("GetCredentials", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceGetCredentialsRequest{
								Project:    testresources.Cluster1().Project,
								Uuid:       testresources.Cluster1().Uuid,
								Expiration: durationpb.New(9 * time.Hour),
							})).Return(connect.NewResponse(&apiv1.ClusterServiceGetCredentialsResponse{
								Kubeconfig: `
apiVersion: v1
kind: Config

clusters:
- name: project--test-external
  cluster:
    server: https://example.com
    insecure-skip-tls-verify: true

users:
- name: myuser-test-external
  user:
    token: fake-token

contexts:
- name: test
  context:
    cluster: project--test-external
    user: myuser-test-external

current-context: test
preferences: {}
`,
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(`
apiVersion: v1
clusters:
- cluster:
    server: https://example.com
  name: test-Some Initiative@metalstack.cloud
contexts:
- context:
    cluster: test-Some Initiative@metalstack.cloud
    user: test-Some Initiative@metalstack.cloud
  name: test-Some Initiative@metalstack.cloud
current-context: test-Some Initiative@metalstack.cloud
kind: Config
users:
- name: test-Some Initiative@metalstack.cloud
  user: {}
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Monitoring(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceGetResponse, *apiv1.Cluster]{
		{
			Name:    "cluster monitoring",
			CmdArgs: []string{"cluster", "monitoring", "--project", testresources.Cluster1().Project, testresources.Cluster1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
								Project: testresources.Cluster1().Project,
								Uuid:    testresources.Cluster1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceGetResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(`
endpoint: endpoint
password: password
username: username
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Reconcile(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceOperateResponse, *apiv1.Cluster]{
		{
			Name:    "cluster monitoring",
			CmdArgs: []string{"cluster", "reconcile", "--project", testresources.Cluster1().Project, testresources.Cluster1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Operate", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceOperateRequest{
								Project: testresources.Cluster1().Project,
								Uuid:    testresources.Cluster1().Uuid,
								Operate: apiv1.Operate_OPERATE_RECONCILE,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceOperateResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Cluster1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ClusterCmd_Status(t *testing.T) {
	tests := []*e2e.Test[apiv1.ClusterServiceGetResponse, *apiv1.Cluster]{
		{
			Name:    "cluster status",
			CmdArgs: []string{"cluster", "status", "--project", testresources.Cluster1().Project, testresources.Cluster1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ClusterServiceGetRequest{
								Project: testresources.Cluster1().Project,
								Uuid:    testresources.Cluster1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ClusterServiceGetResponse{
								Cluster: testresources.Cluster1(),
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(`
   TYPE   MESSAGE                                         REASON           LAST UPDATE                    
✔  Ready  All cluster nodes are reporting healthy status  AllNodesHealthy  Sat, 01 Jan 2000 00:00:00 UTC  

Last Errors:
TIME                           DESCRIPTION  TASK    
Sat, 01 Jan 2000 00:00:00 UTC  failed       someid 
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
