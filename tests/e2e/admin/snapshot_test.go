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

func Test_AdminSnapshotCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.StorageServiceListSnapshotsResponse, []*apiv1.Snapshot]{
		{
			Name:    "list snapshot",
			CmdArgs: []string{"admin", "storage", "snapshot", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Storage: func(m *mock.Mock) {
							m.On("ListSnapshots", mock.Anything, connect.NewRequest(&adminv1.StorageServiceListSnapshotsRequest{})).Return(connect.NewResponse(&adminv1.StorageServiceListSnapshotsResponse{
								Snapshots: []*apiv1.Snapshot{
									testresources.Snapshot1(), testresources.Snapshot2(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            3f4g20b8-406a-4952-b91c-73c38dbed113  prod-db-backup-hour   50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a  
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
		`),
			WantWideTable: new(`
            ID                                    NAME                  SIZE    USAGE   SOURCE VOLUME ID                      SOURCE VOLUME NAME  PROJECT     PARTITION   
            3f4g20b8-406a-4952-b91c-73c38dbed113  prod-db-backup-hour   50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a  
            8f4c20b8-406a-4952-b91c-73c38dbed112  prod-db-backup-daily  50 GiB  40 GiB  1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p  prod-db-volume-01   project-99  us-east-1a
		`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
8f4c20b8-406a-4952-b91c-73c38dbed112 project-99
3f4g20b8-406a-4952-b91c-73c38dbed113 project-99
		`),
			WantMarkdown: new(`
            | ID                                   | NAME                 | SIZE   | USAGE  | SOURCE VOLUME ID                     | SOURCE VOLUME NAME | PROJECT    | PARTITION  |
            |--------------------------------------|----------------------|--------|--------|--------------------------------------|--------------------|------------|------------|
            | 3f4g20b8-406a-4952-b91c-73c38dbed113 | prod-db-backup-hour  | 50 GiB | 40 GiB | 1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p | prod-db-volume-01  | project-99 | us-east-1a |
            | 8f4c20b8-406a-4952-b91c-73c38dbed112 | prod-db-backup-daily | 50 GiB | 40 GiB | 1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p | prod-db-volume-01  | project-99 | us-east-1a |
		`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminSnapshotCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv1.StorageServiceListSnapshotsResponse, *apiv1.Snapshot]{
		{
			Name:    "describe snapshot",
			CmdArgs: []string{"admin", "storage", "snapshot", "describe", testresources.Snapshot1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Storage: func(m *mock.Mock) {
							m.On("ListSnapshots", mock.Anything, connect.NewRequest(&adminv1.StorageServiceListSnapshotsRequest{
								Uuid: &testresources.Snapshot1().Uuid,
							})).Return(connect.NewResponse(&adminv1.StorageServiceListSnapshotsResponse{
								Snapshots: []*apiv1.Snapshot{
									testresources.Snapshot1(),
								},
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
