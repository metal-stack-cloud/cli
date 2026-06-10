package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/stretchr/testify/mock"
)

func Test_AdminClusterCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.ClusterServiceListResponse, []*apiv1.Cluster]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "cluster", "list", "--project", testresources.Cluster2().Project},

			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&adminv1.ClusterServiceListRequest{
								Project: &testresources.Cluster2().Project,
							})).Return(connect.NewResponse(&adminv1.ClusterServiceListResponse{
								Clusters: []*apiv1.Cluster{
									testresources.Cluster2(),
									testresources.Cluster1(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            PROGRESS  ID                                    NAME      PROJECT                               TENANT         PARTITION    VERSION  SIZE   AGE  
            72%       6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack    partition-a  1.25.10  1 - 3  now  
            100%      0c538734-c469-46a0-8efd-98e439d4dc8a  cluster2  c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent.com  partition-b  1.27.9   3 - 6  now
`),
			WantWideTable: new(`
            PROGRESS          ID                                    NAME      PROJECT                               TENANT         PARTITION    VERSION  SIZE   AGE  PURPOSE     OPERATION   API  CONTROL  NODES  SYS  
            72% [Reconcile]   6c631ff1-9038-4ad0-b75e-3ea173b7cdb1  cluster1  c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack    partition-a  1.25.10  1 - 3  now  evaluation  Processing  ✔    ✔        ✗      ✔    
            100% [Reconcile]  0c538734-c469-46a0-8efd-98e439d4dc8a  cluster2  c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent.com  partition-b  1.27.9   3 - 6  now  production  Succeeded   ✔    ✔        ✔      ✔
`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 c40ad996-e1fd-4511-a7bf-418219cb8d95
0c538734-c469-46a0-8efd-98e439d4dc8a c40ad996-e1fd-4511-a7bf-418219cb8d67
			`),
			WantMarkdown: new(`
            | PROGRESS | ID                                   | NAME     | PROJECT                              | TENANT        | PARTITION   | VERSION | SIZE  | AGE |
            |----------|--------------------------------------|----------|--------------------------------------|---------------|-------------|---------|-------|-----|
            | 72%      | 6c631ff1-9038-4ad0-b75e-3ea173b7cdb1 | cluster1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack   | partition-a | 1.25.10 | 1 - 3 | now |
            | 100%     | 0c538734-c469-46a0-8efd-98e439d4dc8a | cluster2 | c40ad996-e1fd-4511-a7bf-418219cb8d67 | x-cellent.com | partition-b | 1.27.9  | 3 - 6 | now |
`)},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminClusterCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv1.ClusterServiceListResponse, *apiv1.Cluster]{
		{
			Name:    "describe",
			CmdArgs: []string{"admin", "cluster", "describe", testresources.Cluster1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Cluster: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&adminv1.ClusterServiceGetRequest{
								Uuid: testresources.Cluster1().Uuid,
							})).Return(connect.NewResponse(&adminv1.ClusterServiceGetResponse{
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

// TODO cluster machine cli
