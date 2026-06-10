package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/stretchr/testify/mock"
)

func Test_VolumeCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceListResponse, []*apiv1.Volume]{
		{
			Name: "list",
			CmdArgs: []string{
				"storage",
				"volume",
				"list",
				"--project",
				testresources.Volume1().Project,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceListRequest{
								Project: testresources.Volume1().Project,
							})).Return(&connect.Response[apiv1.VolumeServiceListResponse]{Msg: &apiv1.VolumeServiceListResponse{
								Volumes: []*apiv1.Volume{
									testresources.Volume1(),
								},
							},
							}, nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
		`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
		`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95
		`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
		`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceGetResponse, *apiv1.Volume]{
		{
			Name: "describe",
			CmdArgs: []string{
				"storage",
				"volume",
				"describe",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
								Uuid: testresources.Volume1().Uuid,
							})).Return(&connect.Response[apiv1.VolumeServiceGetResponse]{Msg: &apiv1.VolumeServiceGetResponse{
								Volume: testresources.Volume1(),
							},
							}, nil)
						},
					},
				}}),
			WantObject: testresources.Volume1(),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceDeleteResponse, *apiv1.Volume]{
		{
			Name: "delete",
			CmdArgs: []string{
				"storage",
				"volume",
				"delete",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceDeleteRequest{
								Uuid: testresources.Volume1().Uuid,
							})).Return(&connect.Response[apiv1.VolumeServiceDeleteResponse]{Msg: &apiv1.VolumeServiceDeleteResponse{
								Volume: testresources.Volume1(),
							},
							}, nil)
						},
					},
				},
			}),
			WantObject: testresources.Volume1(),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceUpdateResponse, *apiv1.Volume]{
		{
			Name: "update",
			CmdArgs: []string{
				"storage",
				"volume",
				"update",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
								Uuid: testresources.Volume1().Uuid,
							})).Return(&connect.Response[apiv1.VolumeServiceUpdateResponse]{Msg: &apiv1.VolumeServiceUpdateResponse{
								Volume: testresources.Volume1(),
							},
							}, nil)
						},
					},
				},
			}),
			WantObject: testresources.Volume1(),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
