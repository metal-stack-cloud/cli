package api_e2e

import (
	"fmt"
	"testing"

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

func Test_TokenCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.TokenServiceListResponse, []*apiv1.Token]{
		{
			Name:    "token list",
			CmdArgs: []string{"token", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.TokenServiceListRequest{})).Return(connect.NewResponse(&apiv1.TokenServiceListResponse{
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

func Test_TokenCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv1.TokenServiceCreateResponse, []*apiv1.Token]{
		{
			Name:    "token apply from file",
			CmdArgs: append([]string{"token", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Token: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TokenServiceCreateRequest{
									Description:  testresources.Token1().Description,
									Permissions:  testresources.Token1().Permissions,
									Expires:      durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
									ProjectRoles: testresources.Token1().ProjectRoles,
									TenantRoles:  testresources.Token1().TenantRoles,
								})).Return(connect.NewResponse(&apiv1.TokenServiceCreateResponse{
									Token:  testresources.Token1(),
									Secret: "secret-01",
								}), nil)
								// FIXME: API does not return a conflict when already exists, so the update functionality does not work!
							},
						},
					},
				}),

			WantDefault: new(`
Make sure to copy your personal access token now as you will not be able to see this again.

secret-01

TYPE            ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
TOKEN_TYPE_API  123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)  
			`),
		},
		{
			Name:    "token apply multiple from file",
			CmdArgs: append([]string{"token", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Token1(), testresources.Token2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Token: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TokenServiceCreateRequest{
									Description:  testresources.Token1().Description,
									Permissions:  testresources.Token1().Permissions,
									Expires:      durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
									ProjectRoles: testresources.Token1().ProjectRoles,
									TenantRoles:  testresources.Token1().TenantRoles,
								})).Return(connect.NewResponse(&apiv1.TokenServiceCreateResponse{
									Token:  testresources.Token1(),
									Secret: "secret-01",
								}), nil)
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TokenServiceCreateRequest{
									Description:  testresources.Token2().Description,
									Permissions:  testresources.Token2().Permissions,
									Expires:      durationpb.New(testresources.Token2().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
									ProjectRoles: testresources.Token2().ProjectRoles,
									TenantRoles:  testresources.Token2().TenantRoles,
								})).Return(connect.NewResponse(&apiv1.TokenServiceCreateResponse{
									Token:  testresources.Token2(),
									Secret: "secret-02",
								}), nil)
								// FIXME: API does not return a conflict when already exists, so the update functionality does not work!
							},
						},
					},
				}),

			WantDefault: new(`
Make sure to copy your personal access token now as you will not be able to see this again.

secret-01

Make sure to copy your personal access token now as you will not be able to see this again.

secret-02

TYPE                ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
TOKEN_TYPE_API      123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)  
TOKEN_TYPE_CONSOLE  ccd042af-7e3a-4458-bc9f-933a26fd783e                     user_Guest      GuestToken                                         1      1      2000-01-02 00:00:00 UTC (in 1d)  
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.TokenServiceUpdateResponse, *apiv1.Token]{
		{
			Name: "token update",
			CmdArgs: []string{"token", "update", testresources.Token1().Uuid,
				"--description", testresources.Token1().Description},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.TokenServiceUpdateRequest{
								Uuid:        testresources.Token1().Uuid,
								Description: &testresources.Token1().Description,
							})).Return(connect.NewResponse(&apiv1.TokenServiceUpdateResponse{
								Token: testresources.Token1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Token1(),
		},
		{
			Name:    "token update from file",
			CmdArgs: append([]string{"token", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Token: func(m *mock.Mock) {
								m.On("Update", mock.Anything, connect.NewRequest(&apiv1.TokenServiceUpdateRequest{
									Uuid:         testresources.Token2().Uuid,
									Description:  &testresources.Token2().Description,
									Permissions:  testresources.Token2().Permissions,
									ProjectRoles: testresources.Token2().ProjectRoles,
									TenantRoles:  testresources.Token2().TenantRoles,
								})).Return(connect.NewResponse(&apiv1.TokenServiceUpdateResponse{
									Token: testresources.Token2(),
								}), nil)
							},
						},
					},
				}),
			WantTable: new(`
            TYPE                ID                                    ADMIN  USER        DESCRIPTION  ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_CONSOLE  ccd042af-7e3a-4458-bc9f-933a26fd783e         user_Guest  GuestToken   1      1      2000-01-02 00:00:00 UTC (in 1d)
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv1.TokenServiceCreateResponse, *apiv1.Token]{
		/* {
			Name: "token create",
			CmdArgs: []string{"token", "create",
				"--description", testresources.Token1().Description,
				"--admin-role", testresources.Token1().AdminRole.String(),
				"--expires", durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())).AsDuration().String()},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TokenServiceCreateRequest{
								Description: testresources.Token1().Description,
								Expires:     durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
								AdminRole:   testresources.Token1().AdminRole,
							})).Return(connect.NewResponse(&apiv1.TokenServiceCreateResponse{
								Token: testresources.Token1(),
							}), nil)
						},
					},
				},
			}),

			WantDefault: new(`

			`),
		}, */
		{
			Name:    "create from file",
			CmdArgs: append([]string{"token", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Token: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TokenServiceCreateRequest{
									Description:  testresources.Token1().Description,
									Expires:      durationpb.New(testresources.Token1().Expires.AsTime().Sub(e2e.TimeBubbleStartTime())),
									Permissions:  testresources.Token1().Permissions,
									ProjectRoles: testresources.Token1().ProjectRoles,
									TenantRoles:  testresources.Token1().TenantRoles,
								})).Return(connect.NewResponse(&apiv1.TokenServiceCreateResponse{
									Token: testresources.Token1(),
								}), nil)
							},
						},
					},
				}),

			WantDefault: new(`
Make sure to copy your personal access token now as you will not be able to see this again.



TYPE            ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
TOKEN_TYPE_API  123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)  
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.TokenServiceRevokeResponse, *apiv1.Token]{
		{
			Name:    "token delete",
			CmdArgs: []string{"token", "delete", testresources.Token1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("Revoke", mock.Anything, connect.NewRequest(&apiv1.TokenServiceRevokeRequest{
								Uuid: testresources.Token1().Uuid,
							})).Return(connect.NewResponse(&apiv1.TokenServiceRevokeResponse{}), nil)
						},
					},
				}}),
			WantDefault: new(fmt.Sprintf("uuid: %s", testresources.Token1().Uuid)),
		},
		{
			Name:    "delete token from file",
			CmdArgs: append([]string{"token", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Token1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Token: func(m *mock.Mock) {
								m.On("Revoke", mock.Anything, connect.NewRequest(&apiv1.TokenServiceRevokeRequest{
									Uuid: testresources.Token1().Uuid,
								})).Return(connect.NewResponse(&apiv1.TokenServiceRevokeResponse{}), nil)
							},
						},
					},
				}),
			Template:     new(`{{ .uuid }}`),
			WantTemplate: new(`123e4567-e89b-12d3-a456-426614174000`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TokenCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.TokenServiceGetResponse, *apiv1.Token]{
		{
			Name:    "token describe",
			CmdArgs: []string{"token", "describe", testresources.Token1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Token: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.TokenServiceGetRequest{
								Uuid: testresources.Token1().Uuid,
							})).Return(connect.NewResponse(&apiv1.TokenServiceGetResponse{
								Token: testresources.Token1(),
							}), nil)
						},
					},
				}}),

			WantObject: testresources.Token1(),
			WantTable: new(`
            TYPE            ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)
`),
			WantWideTable: new(`
            TYPE            ID                                    ADMIN              USER            DESCRIPTION                                        ROLES  PERMS  EXPIRES                          
            TOKEN_TYPE_API  123e4567-e89b-12d3-a456-426614174000  ADMIN_ROLE_EDITOR  user_987654321  Production API read-only token for CI/CD pipeline  2      1      2000-01-02 00:00:00 UTC (in 1d)
`),
			WantMarkdown: new(`
            | TYPE           | ID                                   | ADMIN             | USER           | DESCRIPTION                                       | ROLES | PERMS | EXPIRES                         |
            |----------------|--------------------------------------|-------------------|----------------|---------------------------------------------------|-------|-------|---------------------------------|
            | TOKEN_TYPE_API | 123e4567-e89b-12d3-a456-426614174000 | ADMIN_ROLE_EDITOR | user_987654321 | Production API read-only token for CI/CD pipeline | 2     | 1     | 2000-01-02 00:00:00 UTC (in 1d) |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
