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

func Test_AdminVolumeCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.StorageServiceListVolumesResponse, []*apiv1.Volume]{
		{
			Name: "list",
			CmdArgs: []string{
				"admin",
				"storage",
				"volume",
				"list",
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Storage: func(m *mock.Mock) {
							m.On("ListVolumes", mock.Anything, connect.NewRequest(&adminv1.StorageServiceListVolumesRequest{})).Return(connect.NewResponse(&adminv1.StorageServiceListVolumesResponse{
								Volumes: []*apiv1.Volume{
									testresources.Volume1(),
									testresources.Volume2(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster2      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d67  partition-a  
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
		`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster2      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d67  partition-a         bar=baz  
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
		`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95
0372d029-1077-4e9b-b303-7d64ad5496fd c40ad996-e1fd-4511-a7bf-418219cb8d67
		`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | 0372d029-1077-4e9b-b303-7d64ad5496fd | volume2 | 1.0 KiB | 42 B  | 0        | cluster2     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d67 | partition-a |
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
		`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminVolumeCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv1.StorageServiceListVolumesResponse, *apiv1.Volume]{
		{
			Name: "describe",
			CmdArgs: []string{
				"admin",
				"storage",
				"volume",
				"describe",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Storage: func(m *mock.Mock) {
							m.On("ListVolumes", mock.Anything, connect.NewRequest(&adminv1.StorageServiceListVolumesRequest{
								Uuid: &testresources.Volume1().Uuid,
							})).Return(connect.NewResponse(&adminv1.StorageServiceListVolumesResponse{
								Volumes: []*apiv1.Volume{
									testresources.Volume1(),
								},
							}), nil)
						},
					},
				}}),
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
