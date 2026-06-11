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

func Test_ProjectCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.ProjectServiceListResponse, []*apiv1.Project]{
		{
			Name:    "list projects",
			CmdArgs: []string{"admin", "project", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&adminv1.ProjectServiceListRequest{})).Return(connect.NewResponse(&adminv1.ProjectServiceListResponse{
								Projects: []*apiv1.Project{
									testresources.Project1(),
									testresources.Project2(),
								},
							}), nil)
						},
					},
				},
			}),

			WantTable: new(`
            ID                                    TENANT       NAME                      DESCRIPTION                                                           CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative           Internal research and development for something.                      2000-01-01 00:00:00 UTC  
            c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent    Some Initiative Number 2  Internal research and development for something even more important.  2000-01-01 00:00:00 UTC
`),
			WantWideTable: new(`
            ID                                    TENANT       NAME                      DESCRIPTION                                                           CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative           Internal research and development for something.                      2000-01-01 00:00:00 UTC  
            c40ad996-e1fd-4511-a7bf-418219cb8d67  x-cellent    Some Initiative Number 2  Internal research and development for something even more important.  2000-01-01 00:00:00 UTC
`),
			Template: new("{{ .tenant }} {{ .name }}"),
			WantTemplate: new(`
metal-stack Some Initiative
x-cellent Some Initiative Number 2
			`),
			WantMarkdown: new(`
            | ID                                   | TENANT      | NAME                     | DESCRIPTION                                                          | CREATION DATE           |
            |--------------------------------------|-------------|--------------------------|----------------------------------------------------------------------|-------------------------|
            | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | Some Initiative          | Internal research and development for something.                     | 2000-01-01 00:00:00 UTC |
            | c40ad996-e1fd-4511-a7bf-418219cb8d67 | x-cellent   | Some Initiative Number 2 | Internal research and development for something even more important. | 2000-01-01 00:00:00 UTC |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
