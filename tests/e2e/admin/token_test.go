package admin_e2e

import (
	"fmt"
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

func Test_AdminTokenCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.TokenServiceListResponse, []*apiv1.Token]{
		{
			Name:    "token list",
			CmdArgs: []string{"admin", "token", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&adminv1.TokenServiceListRequest{})).Return(connect.NewResponse(&adminv1.TokenServiceListResponse{
								Tokens: []*apiv1.Token{
									testresources.Token1(), testresources.Token2(),
								},
							}), nil)
						},
					},
				},
			}),

			WantTable: new(`
            TYPE                ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API      123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)  
            TOKEN_TYPE_CONSOLE  ccd042af-7e3a-4458-bc9f-933a26fd783e                     user_Guest      GuestToken                                         1      1      2000-01-02 00:00:00 UTC (in 1d)
`),
			WantWideTable: new(`
            TYPE                ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API      123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)  
            TOKEN_TYPE_CONSOLE  ccd042af-7e3a-4458-bc9f-933a26fd783e                     user_Guest      GuestToken                                         1      1      2000-01-02 00:00:00 UTC (in 1d)
`),
			Template: new("{{ .user_id }} {{ .token_type }}"),
			WantTemplate: new(`
user_987654321 1
user_Guest 2
			`),
			WantMarkdown: new(`
            | TYPE               | ID                                   | ADMIN             | USER           | DESCRIPTION                                       | ROLES | PERMS | EXPIRES                         |
            |--------------------|--------------------------------------|-------------------|----------------|---------------------------------------------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API     | 123e4567-e89b-12d3-a456-426614174000 | ADMIN_ROLE_EDITOR | user_987654321 | Production API read-only token for CI/CD pipeline | 2     | 1     | 2000-01-02 00:00:00 UTC (in 1d) |
            | TOKEN_TYPE_CONSOLE | ccd042af-7e3a-4458-bc9f-933a26fd783e |                   | user_Guest     | GuestToken                                        | 1     | 1     | 2000-01-02 00:00:00 UTC (in 1d) |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTokenCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[adminv1.TokenServiceRevokeResponse, *apiv1.Token]{
		{
			Name:    "token delete",
			CmdArgs: []string{"admin", "token", "delete", testresources.Token1().Uuid, "--user", testresources.Token1().UserId},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("Revoke", mock.Anything, connect.NewRequest(&adminv1.TokenServiceRevokeRequest{
								Uuid:   testresources.Token1().Uuid,
								UserId: testresources.Token1().UserId,
							})).Return(connect.NewResponse(&adminv1.TokenServiceRevokeResponse{}), nil)
						},
					},
				}}),
			WantDefault: new(fmt.Sprintf("uuid: %s", testresources.Token1().Uuid)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
