package api_e2e

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/dustin/go-humanize"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_TenantCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceGetResponse, *apiv1.Tenant]{
		{
			Name:    "get tenant",
			CmdArgs: []string{"tenant", "describe", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.TenantServiceGetRequest{
								Login: testresources.Tenant1().Login,
							})).Return(connect.NewResponse(&apiv1.TenantServiceGetResponse{
								Tenant: testresources.Tenant1(),
								TenantMembers: []*apiv1.TenantMember{
									testresources.TenantMember1(),
								},
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Tenant1(),
			WantTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
`),
			WantWideTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
`),
			Template: new("{{ .login }} {{ .name }} {{ .email }} {{ .oauth_provider }} {{ if .admitted }}true{{ else }}false{{ end }} {{ if .terms_and_conditions.accepted }}true{{ else }}false{{ end }}"),
			WantTemplate: new(`
metal-stack metal-stack metal-stack@mail.com 1 true true
			`),
			WantMarkdown: new(`
            | ID          | NAME        | EMAIL                | REGISTERED | PROVIDER               | ADMITTED | TERMS AND CONDITIONS |
            |-------------|-------------|----------------------|------------|------------------------|----------|----------------------|
            | metal-stack | metal-stack | metal-stack@mail.com | now        | O_AUTH_PROVIDER_GITHUB | true     | true                 |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceListResponse, []*apiv1.Tenant]{
		{
			Name:    "list tenants",
			CmdArgs: []string{"tenant", "list", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.TenantServiceListRequest{})).Return(connect.NewResponse(&apiv1.TenantServiceListResponse{
								Tenants: []*apiv1.Tenant{
									testresources.Tenant1(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
`),
			WantWideTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
`),
			Template: new("{{ .login }} {{ .name }} {{ .email }} {{ .oauth_provider }} {{ if .admitted }}true{{ else }}false{{ end }} {{ if .terms_and_conditions.accepted }}true{{ else }}false{{ end }}"),
			WantTemplate: new(`
metal-stack metal-stack metal-stack@mail.com 1 true true
			`),
			WantMarkdown: new(`
            | ID          | NAME        | EMAIL                | REGISTERED | PROVIDER               | ADMITTED | TERMS AND CONDITIONS |
            |-------------|-------------|----------------------|------------|------------------------|----------|----------------------|
            | metal-stack | metal-stack | metal-stack@mail.com | now        | O_AUTH_PROVIDER_GITHUB | true     | true                 |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceCreateResponse, *apiv1.Tenant]{
		{
			Name:    "create",
			CmdArgs: []string{"tenant", "create", "--name", testresources.Tenant1().Name, "--description", testresources.Tenant1().Description, "--email", testresources.Tenant1().Email},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
								Name:        testresources.Tenant1().Name,
								Description: &testresources.Tenant1().Description,
								Email:       &testresources.Tenant1().Email,
							})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantObject:      testresources.Tenant1(),
			WantProtoObject: testresources.Tenant1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"tenant", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
								Name:        testresources.Tenant1().Name,
								Description: &testresources.Tenant1().Description,
								Email:       &testresources.Tenant1().Email,
								AvatarUrl:   &testresources.Tenant1().AvatarUrl,
								PhoneNumber: &testresources.Tenant1().PhoneNumber,
							})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
					`),
		},
		{
			Name:    "create many from file",
			CmdArgs: append([]string{"tenant", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Tenant1(), testresources.Tenant2()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
								Name:        testresources.Tenant2().Name,
								Description: &testresources.Tenant2().Description,
								Email:       &testresources.Tenant2().Email,
								AvatarUrl:   &testresources.Tenant2().AvatarUrl,
								PhoneNumber: &testresources.Tenant2().PhoneNumber,
							})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
								Tenant: testresources.Tenant2(),
							}), nil)
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
								Name:        testresources.Tenant1().Name,
								Description: &testresources.Tenant1().Description,
								Email:       &testresources.Tenant1().Email,
								AvatarUrl:   &testresources.Tenant1().AvatarUrl,
								PhoneNumber: &testresources.Tenant1().PhoneNumber,
							})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID             NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack    metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true                  
            x-cellent.com  x-cellent    x-cellent@mail.com    now         O_AUTH_PROVIDER_GITHUB  false     false
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceDeleteResponse, *apiv1.Tenant]{
		{
			Name:    "delete",
			CmdArgs: []string{"tenant", "delete", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.TenantServiceDeleteRequest{
								Login: testresources.Tenant1().Login,
							})).Return(connect.NewResponse(&apiv1.TenantServiceDeleteResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Tenant1(),
		},
		{
			Name:    "delete from file",
			CmdArgs: append([]string{"tenant", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.TenantServiceDeleteRequest{
								Login: testresources.Tenant1().Login,
							})).Return(connect.NewResponse(&apiv1.TenantServiceDeleteResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceUpdateResponse, *apiv1.Tenant]{
		{
			Name:    "update",
			CmdArgs: []string{"tenant", "update", testresources.Tenant1().Login, "--name", "new-name", "--description", "new-desc"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						User: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.UserServiceGetRequest{})).Return(connect.NewResponse(&apiv1.UserServiceGetResponse{
								User: testresources.User1(),
							}), nil)
						},
						Tenant: func(m *mock.Mock) {
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.TenantServiceUpdateRequest{
								Login:       testresources.Tenant1().Login,
								Name:        new("new-name"),
								Description: new("new-desc"),
							})).Return(connect.NewResponse(&apiv1.TenantServiceUpdateResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Tenant1(),
			WantTable: new(`
			ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
			`),
		},
		{
			Name:    "update from file",
			CmdArgs: append([]string{"tenant", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant1()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.TenantServiceUpdateRequest{
								Login: testresources.Tenant1().Login,
								Name:  &testresources.Tenant1().Name,
								Email: &testresources.Tenant1().Email,
							})).Return(connect.NewResponse(&apiv1.TenantServiceUpdateResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID           NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack  metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true
					`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceUpdateResponse, *apiv1.Tenant]{
		{
			Name:    "apply from file",
			CmdArgs: append([]string{"tenant", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Tenant: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
									Name:        testresources.Tenant2().Name,
									Description: &testresources.Tenant2().Description,
									Email:       &testresources.Tenant2().Email,
									AvatarUrl:   &testresources.Tenant2().AvatarUrl,
									PhoneNumber: &testresources.Tenant2().PhoneNumber,
								})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
									Tenant: testresources.Tenant2(),
								}), nil)
							},
						},
					},
				},
			),
			WantTable: new(`
            ID             NAME       EMAIL               REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            x-cellent.com  x-cellent  x-cellent@mail.com  now         O_AUTH_PROVIDER_GITHUB  false     false
			`),
		},
		{
			Name:    "apply many from file",
			CmdArgs: append([]string{"tenant", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Tenant1(), testresources.Tenant2()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
								Name:        testresources.Tenant2().Name,
								Description: &testresources.Tenant2().Description,
								Email:       &testresources.Tenant2().Email,
								AvatarUrl:   &testresources.Tenant2().AvatarUrl,
								PhoneNumber: &testresources.Tenant2().PhoneNumber,
							})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
								Tenant: testresources.Tenant2(),
							}), nil)
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
								Name:        testresources.Tenant1().Name,
								Description: &testresources.Tenant1().Description,
								Email:       &testresources.Tenant1().Email,
								AvatarUrl:   &testresources.Tenant1().AvatarUrl,
								PhoneNumber: &testresources.Tenant1().PhoneNumber,
							})).Return(connect.NewResponse(&apiv1.TenantServiceCreateResponse{
								Tenant: testresources.Tenant1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID             NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack    metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true                  
            x-cellent.com  x-cellent    x-cellent@mail.com    now         O_AUTH_PROVIDER_GITHUB  false     false
					`),
		},
		{
			Name:    "apply already exists",
			CmdArgs: append([]string{"tenant", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Tenant2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Tenant: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.TenantServiceCreateRequest{
									Name:        testresources.Tenant2().Name,
									Description: &testresources.Tenant2().Description,
									Email:       &testresources.Tenant2().Email,
									AvatarUrl:   &testresources.Tenant2().AvatarUrl,
									PhoneNumber: &testresources.Tenant2().PhoneNumber,
								})).Return(nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")))
							},
						},
					},
				},
			),
			WantErr: fmt.Errorf("error creating entity: failed to create tenant: already_exists: already exists"),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_ListMembers(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceGetResponse, []apiv1.TenantMember]{
		{
			Name:    "list tenant members",
			CmdArgs: []string{"tenant", "member", "list", "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.TenantServiceGetRequest{
								Login: testresources.Tenant1().Login,
							})).Return(connect.NewResponse(&apiv1.TenantServiceGetResponse{
								Tenant:        testresources.Tenant2(),
								TenantMembers: []*apiv1.TenantMember{testresources.TenantMember1(), testresources.TenantMember2()},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID             ROLE               SINCE  
            metal-stack    TENANT_ROLE_OWNER  now    
            x-cellent.com  TENANT_ROLE_OWNER  now
			`),
			WantWideTable: new(`
            ID             ROLE               SINCE  
            metal-stack    TENANT_ROLE_OWNER  now    
            x-cellent.com  TENANT_ROLE_OWNER  now
			`),
			Template: new("{{ .id }} {{ .role }}"),
			WantTemplate: new(`
metal-stack 1
x-cellent.com 1
			`),
			WantMarkdown: new(`
            | ID            | ROLE              | SINCE |
            |---------------|-------------------|-------|
            | metal-stack   | TENANT_ROLE_OWNER | now   |
            | x-cellent.com | TENANT_ROLE_OWNER | now   |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_DeleteMember(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceRemoveMemberResponse, string]{
		{
			Name:    "delete tenant member",
			CmdArgs: []string{"tenant", "member", "remove", testresources.TenantMember1().Id, "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("RemoveMember", mock.Anything, connect.NewRequest(&apiv1.TenantServiceRemoveMemberRequest{
								Login:    testresources.Tenant1().Login,
								MemberId: testresources.TenantMember1().Id,
							})).Return(connect.NewResponse(&apiv1.TenantServiceRemoveMemberResponse{}), nil)
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("✔ successfully removed member \"%s\"", testresources.TenantMember1().Id)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_UpdateMember(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceUpdateMemberResponse, *apiv1.TenantMember]{
		{
			Name:    "update tenant member",
			CmdArgs: []string{"tenant", "member", "update", testresources.TenantMember1().Id, "--tenant", testresources.Tenant1().Login, "--role", testresources.TenantMember1().Role.String()},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("UpdateMember", mock.Anything, connect.NewRequest(&apiv1.TenantServiceUpdateMemberRequest{
								Login:    testresources.Tenant1().Login,
								MemberId: testresources.TenantMember1().Id,
								Role:     testresources.TenantMember1().Role,
							})).Return(connect.NewResponse(&apiv1.TenantServiceUpdateMemberResponse{
								TenantMember: testresources.TenantMember1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.TenantMember1(),
			WantTable: new(`
            ID           ROLE               SINCE  
            metal-stack  TENANT_ROLE_OWNER  now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_ListInvites(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceInvitesListResponse, apiv1.TenantInvite]{
		{
			Name:    "list invites",
			CmdArgs: []string{"tenant", "invite", "list", "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("InvitesList", mock.Anything, connect.NewRequest(&apiv1.TenantServiceInvitesListRequest{
								Login: testresources.Tenant1().Login,
							})).Return(connect.NewResponse(&apiv1.TenantServiceInvitesListResponse{
								Invites: []*apiv1.TenantInvite{testresources.TenantInvite1(), testresources.TenantInvite2()},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            SECRET  TENANT       INVITED BY   ROLE                EXPIRES IN       
            secret  acme-corp    acme-corp    TENANT_ROLE_EDITOR  2 days from now  
            secret  metal-stack  metal-stack  TENANT_ROLE_VIEWER  2 days from now
			`),
			WantWideTable: new(`
            SECRET  TENANT       INVITED BY   ROLE                EXPIRES IN       
            secret  acme-corp    acme-corp    TENANT_ROLE_EDITOR  2 days from now  
            secret  metal-stack  metal-stack  TENANT_ROLE_VIEWER  2 days from now
			`),
			Template: new("{{ .tenant }} {{ .role }}"),
			WantTemplate: new(`
acme-corp 2
metal-stack 3
			`),
			WantMarkdown: new(`
            | SECRET | TENANT      | INVITED BY  | ROLE               | EXPIRES IN      |
            |--------|-------------|-------------|--------------------|-----------------|
            | secret | acme-corp   | acme-corp   | TENANT_ROLE_EDITOR | 2 days from now |
            | secret | metal-stack | metal-stack | TENANT_ROLE_VIEWER | 2 days from now |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_DeleteInvite(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceInviteDeleteResponse, string]{
		{
			Name:    "delete invite",
			CmdArgs: []string{"tenant", "invite", "delete", testresources.TenantInvite1().Secret, "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("InviteDelete", mock.Anything, connect.NewRequest(&apiv1.TenantServiceInviteDeleteRequest{
								Login:  testresources.Tenant1().Login,
								Secret: testresources.TenantInvite1().Secret,
							})).Return(connect.NewResponse(&apiv1.TenantServiceInviteDeleteResponse{}), nil)
						},
					},
				},
			}),
			WantMarkdown: new(""),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_CreateInvite(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceInviteResponse, string]{
		{
			Name:    "create invite",
			CmdArgs: []string{"tenant", "invite", "generate-join-secret", "--role", testresources.TenantInvite1().Role.String(), "--tenant", testresources.Tenant1().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Asset: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.AssetServiceListRequest{})).Return(connect.NewResponse(&apiv1.AssetServiceListResponse{
								Assets: []*apiv1.Asset{},
							}), nil)
						},
						Tenant: func(m *mock.Mock) {
							m.On("Invite", mock.Anything, connect.NewRequest(&apiv1.TenantServiceInviteRequest{
								Login: testresources.Tenant1().Login,
								Role:  testresources.TenantInvite1().Role,
							})).Return(connect.NewResponse(&apiv1.TenantServiceInviteResponse{
								Invite: testresources.TenantInvite1(),
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("You can share this secret with the member to join, it expires in %s:\n\n%s (https://console.metalstack.cloud/organization-invite/%s)",
				humanize.RelTime(e2e.TimeBubbleStartTime(), testresources.ProjectInvite1().ExpiresAt.AsTime(), "from now", "ago"),
				testresources.ProjectInvite1().Secret,
				testresources.ProjectInvite1().Secret,
			)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_Join(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceInviteAcceptResponse, string]{
		{
			Name:    "join",
			CmdArgs: []string{"tenant", "invite", "join", testresources.TenantInvite1().Secret},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("InviteGet", mock.Anything, connect.NewRequest(&apiv1.TenantServiceInviteGetRequest{
								Secret: testresources.TenantInvite1().Secret,
							})).Return(connect.NewResponse(&apiv1.TenantServiceInviteGetResponse{
								Invite: testresources.TenantInvite1(),
							}), nil)
							m.On("InviteAccept", mock.Anything, connect.NewRequest(&apiv1.TenantServiceInviteAcceptRequest{
								Secret: testresources.TenantInvite1().Secret,
							})).Return(connect.NewResponse(&apiv1.TenantServiceInviteAcceptResponse{
								Tenant:     testresources.TenantInvite1().TargetTenant,
								TenantName: testresources.TenantInvite1().TargetTenantName,
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("Do you want to join tenant \"%s\" as %s? [Y/n] ✔ successfully joined tenant \"%s\"",
				testresources.TenantInvite1().TenantName,
				testresources.TenantInvite1().Role.String(),
				testresources.TenantInvite1().TenantName)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_TenantCmd_RequestAdmission(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceRequestAdmissionResponse, string]{
		{
			Name:    "request admission",
			CmdArgs: []string{"tenant", "request-admission", testresources.Tenant1().Login, testresources.User1().Email},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						User: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.UserServiceGetRequest{})).Return(connect.NewResponse(&apiv1.UserServiceGetResponse{
								User: testresources.User1(),
							}), nil).Maybe()
						},
						Asset: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.AssetServiceListRequest{})).Return(connect.NewResponse(&apiv1.AssetServiceListResponse{
								Environment: testresources.Environment1(),
							}), nil)
						},
						Tenant: func(m *mock.Mock) {
							m.On("RequestAdmission", mock.Anything, connect.NewRequest(&apiv1.TenantServiceRequestAdmissionRequest{
								Login:                      testresources.Tenant1().Login,
								Email:                      testresources.User1().Email,
								Name:                       testresources.Tenant1().Name,
								EmailConsent:               false,
								AcceptedTermsAndConditions: true,
							})).Return(connect.NewResponse(&apiv1.TenantServiceRequestAdmissionResponse{}), nil)
						},
					},
				},
			}),
			WantDefault: new(fmt.Sprintf("The terms and conditions can be found on %s. Do you accept? [Y/n] Your admission request has been submitted. We will contact you as soon as possible.",
				*testresources.Environment1().TermsAndConditionsUrl)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
