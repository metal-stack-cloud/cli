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

func Test_SnapshotCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.SnapshotServiceListResponse, []*apiv1.Snapshot]{
		{
			Name: "list snapshot",
			CmdArgs: []string{
				"storage",
				"snapshot",
				"list",
				"--project",
				testresources.Snapshot1().Project,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Snapshot: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.SnapshotServiceListRequest{
								Project: testresources.Snapshot1().Project,
							})).Return(connect.NewResponse(&apiv1.SnapshotServiceListResponse{
								Snapshots: []*apiv1.Snapshot{
									testresources.Snapshot1(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
		`),
			WantWideTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
		`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
8f4c20b8-406a-4952-b91c-73c38dbed112 project-99
		`),
			WantMarkdown: new(`
            | ID                                   | NAME                 | SIZE   | USAGE  | SOURCE VOLUME ID                     | SOURCE VOLUME NAME | PROJECT    | PARTITION  |
            |--------------------------------------|----------------------|--------|--------|--------------------------------------|--------------------|------------|------------|
            | 8f4c20b8-406a-4952-b91c-73c38dbed112 | prod-db-backup-daily | 50 GiB | 40 GiB | 1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p | prod-db-volume-01  | project-99 | us-east-1a |
		`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_SnapshotCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.SnapshotServiceGetResponse, *apiv1.Snapshot]{
		{
			Name: "describe snapshot",
			CmdArgs: []string{
				"storage",
				"snapshot",
				"describe",
				testresources.Snapshot1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Snapshot: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.SnapshotServiceGetRequest{
								Uuid: testresources.Snapshot1().Uuid,
							})).Return(connect.NewResponse(&apiv1.SnapshotServiceGetResponse{
								Snapshot: testresources.Snapshot1(),
							}), nil)
						},
					},
				}}),
			WantObject: testresources.Snapshot1(),
			WantTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
`),
			WantWideTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`8f4c20b8-406a-4952-b91c-73c38dbed112 project-99`),
			WantMarkdown: new(`
            | ID                                   | NAME                 | SIZE   | USAGE  | SOURCE VOLUME ID                     | SOURCE VOLUME NAME | PROJECT    | PARTITION  |
            |--------------------------------------|----------------------|--------|--------|--------------------------------------|--------------------|------------|------------|
            | 8f4c20b8-406a-4952-b91c-73c38dbed112 | prod-db-backup-daily | 50 GiB | 40 GiB | 1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p | prod-db-volume-01  | project-99 | us-east-1a |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_SnapshotCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.SnapshotServiceDeleteResponse, *apiv1.Snapshot]{
		{
			Name: "delete",
			CmdArgs: []string{
				"storage",
				"snapshot",
				"delete",
				testresources.Snapshot1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Snapshot: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.SnapshotServiceDeleteRequest{
								Uuid: testresources.Snapshot1().Uuid,
							})).Return(connect.NewResponse(&apiv1.SnapshotServiceDeleteResponse{
								Snapshot: testresources.Snapshot1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Snapshot1(),
			WantTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
`),
			WantWideTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`8f4c20b8-406a-4952-b91c-73c38dbed112 project-99`),
			WantMarkdown: new(`
            | ID                                   | NAME                 | SIZE   | USAGE  | SOURCE VOLUME ID                     | SOURCE VOLUME NAME | PROJECT    | PARTITION  |
            |--------------------------------------|----------------------|--------|--------|--------------------------------------|--------------------|------------|------------|
            | 8f4c20b8-406a-4952-b91c-73c38dbed112 | prod-db-backup-daily | 50 GiB | 40 GiB | 1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p | prod-db-volume-01  | project-99 | us-east-1a |
`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
