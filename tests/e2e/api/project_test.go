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

func Test_ProjectCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceListResponse, []*apiv1.Project]{
		{
			Name:    "list projects",
			CmdArgs: []string{"project", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceListRequest{})).Return(connect.NewResponse(&apiv1.ProjectServiceListResponse{
								Projects: []*apiv1.Project{
									testresources.Project1(),
								},
							}), nil)
						},
					},
				},
			}),

			WantTable: new(`
            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
`),
			WantWideTable: new(`
            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
`),
			Template:     new("{{ .tenant }} {{ .name }}"),
			WantTemplate: new(`metal-stack Some Initiative`),
			WantMarkdown: new(`
            | ID                                   | TENANT      | NAME            | DESCRIPTION                                      | CREATION DATE           |
            |--------------------------------------|-------------|-----------------|--------------------------------------------------|-------------------------|
            | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | Some Initiative | Internal research and development for something. | 2000-01-01 00:00:00 UTC |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceCreateResponse, []*apiv1.Project]{
		{
			Name:    "apply project",
			CmdArgs: append([]string{"project", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Project1(), testresources.Project2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Project: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceCreateRequest{
									Login:       testresources.Project2().Tenant,
									Name:        testresources.Project2().Name,
									Description: testresources.Project2().Description,
								})).Return(connect.NewResponse(&apiv1.ProjectServiceCreateResponse{
									Project: testresources.Project2(),
								}), nil)
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceCreateRequest{
									Login:       testresources.Project1().Tenant,
									Name:        testresources.Project1().Name,
									Description: testresources.Project1().Description,
								})).Return(connect.NewResponse(&apiv1.ProjectServiceCreateResponse{
									Project: testresources.Project1(),
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
			WantMarkdown: new(`
            | ID                                   | TENANT      | NAME                     | DESCRIPTION                                                          | CREATION DATE           |
            |--------------------------------------|-------------|--------------------------|----------------------------------------------------------------------|-------------------------|
            | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | Some Initiative          | Internal research and development for something.                     | 2000-01-01 00:00:00 UTC |
            | c40ad996-e1fd-4511-a7bf-418219cb8d67 | x-cellent   | Some Initiative Number 2 | Internal research and development for something even more important. | 2000-01-01 00:00:00 UTC |
			`),
		},
		{
			Name:    "apply exits",
			CmdArgs: append([]string{"project", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Project: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceCreateRequest{
									Login:       testresources.Project1().Tenant,
									Name:        testresources.Project1().Name,
									Description: testresources.Project1().Description,
								})).Return(nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("already exists")))
							},
						},
					},
				}),

			WantErr: fmt.Errorf("error creating entity: failed to create project: already_exists: already exists"),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceUpdateResponse, *apiv1.Project]{
		{
			Name: "update project",
			CmdArgs: []string{"project", "update", testresources.Project1().Uuid,
				"--description", testresources.Project1().Description,
				"--name", testresources.Project1().Name},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceUpdateRequest{
								Project:     testresources.Project1().Uuid,
								Name:        &testresources.Project1().Name,
								Description: &testresources.Project1().Description,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceUpdateResponse{
								Project: testresources.Project1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Project1(),
			WantDefault: new(`
avatarUrl: https://cdn.example.com/avatars/me.png
createdAt: "2000-01-01T00:00:00Z"
description: Internal research and development for something.
name: Some Initiative
tenant: metal-stack
uuid: c40ad996-e1fd-4511-a7bf-418219cb8d95
			`),
		},
		{
			Name:    "update project from file",
			CmdArgs: append([]string{"project", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Project: func(m *mock.Mock) {
								m.On("Update", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceUpdateRequest{
									Project:     testresources.Project1().Uuid,
									Name:        &testresources.Project1().Name,
									Description: &testresources.Project1().Description,
								})).Return(connect.NewResponse(&apiv1.ProjectServiceUpdateResponse{
									Project: testresources.Project1(),
								}), nil)
							},
						},
					},
				}),
			WantTable: new(`
            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceCreateResponse, *apiv1.Project]{
		{
			Name: "create project",
			CmdArgs: []string{"project", "create",
				"--description", testresources.Project1().Description,
				"--name", testresources.Project1().Name,
				"--tenant", testresources.Tenant1().Login,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceCreateRequest{
								Login:       testresources.Tenant1().Login,
								Name:        testresources.Project1().Name,
								Description: testresources.Project1().Description,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceCreateResponse{
								Project: testresources.Project1(),
							}), nil)
						},
					},
				},
			}),

			WantObject: testresources.Project1(),
		},
		{
			Name:    "create project from file",
			CmdArgs: append([]string{"project", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Project: func(m *mock.Mock) {
								m.On("Create", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceCreateRequest{
									Login:       testresources.Tenant1().Login,
									Name:        testresources.Project1().Name,
									Description: testresources.Project1().Description,
								})).Return(connect.NewResponse(&apiv1.ProjectServiceCreateResponse{
									Project: testresources.Project1(),
								}), nil)
							},
						},
					},
				}),

			WantTable: new(`
			            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE
			            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
								`),
			WantWideTable: new(`
			            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE
			            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
								`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceDeleteResponse, *apiv1.Project]{
		{
			Name:    "delete project",
			CmdArgs: []string{"project", "delete", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceDeleteRequest{
								Project: testresources.Project1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceDeleteResponse{
								Project: testresources.Project1(),
							}), nil)
						},
					},
				}}),

			WantObject: testresources.Project1(),
		},
		{
			Name:    "delete project from file",
			CmdArgs: append([]string{"project", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Project1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							Project: func(m *mock.Mock) {
								m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceDeleteRequest{
									Project: testresources.Project1().Uuid,
								})).Return(connect.NewResponse(&apiv1.ProjectServiceDeleteResponse{
									Project: testresources.Project1(),
								}), nil)
							},
						},
					},
				}),

			WantDefault: new(`
ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE            
c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC  
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceGetResponse, *apiv1.Project]{
		{
			Name:    "describe project",
			CmdArgs: []string{"project", "describe", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
								Project: testresources.Project1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceGetResponse{
								Project: testresources.Project1(),
							}), nil)
						},
					},
				}}),
			WantObject: testresources.Project1(),
			WantTable: new(`
            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
`),
			WantWideTable: new(`
            ID                                    TENANT       NAME             DESCRIPTION                                       CREATION DATE            
            c40ad996-e1fd-4511-a7bf-418219cb8d95  metal-stack  Some Initiative  Internal research and development for something.  2000-01-01 00:00:00 UTC
`),
			Template:     new("{{ .tenant }} {{ .name }}"),
			WantTemplate: new(`metal-stack Some Initiative`),
			WantMarkdown: new(`
            | ID                                   | TENANT      | NAME            | DESCRIPTION                                      | CREATION DATE           |
            |--------------------------------------|-------------|-----------------|--------------------------------------------------|-------------------------|
            | c40ad996-e1fd-4511-a7bf-418219cb8d95 | metal-stack | Some Initiative | Internal research and development for something. | 2000-01-01 00:00:00 UTC |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_ListInvites(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceInvitesListResponse, apiv1.ProjectInvite]{
		{
			Name:    "list invites",
			CmdArgs: []string{"project", "invite", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("InvitesList", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceInvitesListRequest{})).Return(connect.NewResponse(&apiv1.ProjectServiceInvitesListResponse{
								Invites: []*apiv1.ProjectInvite{testresources.ProjectInvite1(), testresources.ProjectInvite2()},
							}), nil)
						},
					},
				}}),
			WantTable: new(`
            SECRET  PROJECT                               ROLE                 EXPIRES IN       
            secret  0d81bca7-73f6-4da3-8397-4a8c52a0c583  PROJECT_ROLE_EDITOR  2 days from now  
            secret  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  PROJECT_ROLE_EDITOR  2 days from now
			`),
			WantWideTable: new(`
            SECRET  PROJECT                               ROLE                 EXPIRES IN       
            secret  0d81bca7-73f6-4da3-8397-4a8c52a0c583  PROJECT_ROLE_EDITOR  2 days from now  
            secret  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  PROJECT_ROLE_EDITOR  2 days from now
			`),
			Template: new("{{ .project }} {{ .role }}"),
			WantTemplate: new(`
0d81bca7-73f6-4da3-8397-4a8c52a0c583 2
f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c 2
			`),
			WantMarkdown: new(`
            | SECRET | PROJECT                              | ROLE                | EXPIRES IN      |
            |--------|--------------------------------------|---------------------|-----------------|
            | secret | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | PROJECT_ROLE_EDITOR | 2 days from now |
            | secret | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | PROJECT_ROLE_EDITOR | 2 days from now |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_DeleteInvite(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceInviteDeleteResponse, string]{
		{
			Name:    "delete",
			CmdArgs: []string{"project", "invite", "delete", testresources.ProjectInvite1().Secret, "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("InviteDelete", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceInviteDeleteRequest{
								Project: testresources.Project1().Uuid,
								Secret:  testresources.ProjectInvite1().Secret,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceInviteDeleteResponse{}), nil)
						},
					},
				}}),
			WantDefault: new(""),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_CreateInvite(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceInviteResponse, string]{
		{
			Name:    "create invite",
			CmdArgs: []string{"project", "invite", "generate-join-secret", "--role", testresources.ProjectInvite1().Role.String(), "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Asset: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.AssetServiceListRequest{})).Return(connect.NewResponse(&apiv1.AssetServiceListResponse{
								Assets: []*apiv1.Asset{},
							}), nil)
						},
						Project: func(m *mock.Mock) {
							m.On("Invite", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceInviteRequest{
								Project: testresources.Project1().Uuid,
								Role:    testresources.ProjectInvite1().Role,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceInviteResponse{
								Invite: testresources.ProjectInvite1(),
							}), nil)
						},
					},
				}}),
			WantDefault: new(fmt.Sprintf("You can share this secret with the member to join, it expires in %s:\n\n%s (https://console.metalstack.cloud/project-invite/%s)",
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

func Test_ProjectCmd_Join(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceInviteAcceptResponse, string]{
		{
			Name:    "join",
			CmdArgs: []string{"project", "invite", "join", testresources.ProjectInvite1().Secret},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("InviteGet", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceInviteGetRequest{
								Secret: testresources.ProjectInvite1().Secret,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceInviteGetResponse{
								Invite: testresources.ProjectInvite1(),
							}), nil)
							m.On("InviteAccept", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceInviteAcceptRequest{
								Secret: testresources.ProjectInvite1().Secret,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceInviteAcceptResponse{
								Project:     testresources.ProjectInvite1().Project,
								ProjectName: testresources.ProjectInvite1().ProjectName,
							}), nil)
						},
					},
				}}),
			WantDefault: new(fmt.Sprintf("Do you want to join project \"%s\" as %s? [Y/n] ✔ successfully joined project \"%s\"",
				testresources.ProjectInvite1().ProjectName,
				testresources.ProjectInvite1().Role.String(),
				testresources.ProjectInvite1().ProjectName)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_ListMembers(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceGetResponse, []apiv1.ProjectMember]{
		{
			Name:    "list project members",
			CmdArgs: []string{"project", "member", "list", "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceGetRequest{
								Project: testresources.Project1().Uuid,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceGetResponse{
								Project:        testresources.Project1(),
								ProjectMembers: []*apiv1.ProjectMember{testresources.ProjectMember1(), testresources.ProjectMember2()},
							}), nil)
						},
					},
				}}),
			WantTable: new(`
            ID                                    ROLE                 INHERITED  SINCE  
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER   false      now    
            40c0da4b-9eb9-4371-91aa-1ae62193fa54  PROJECT_ROLE_EDITOR  true       now    
			`),
			WantWideTable: new(`
            ID                                    ROLE                 INHERITED  SINCE  
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER   false      now    
            40c0da4b-9eb9-4371-91aa-1ae62193fa54  PROJECT_ROLE_EDITOR  true       now
			`),
			Template: new("{{ .id }} {{ .role }}"),
			WantTemplate: new(`
16d6e8ba-f574-494f-8d5e-74f6cb2d8db0 1
40c0da4b-9eb9-4371-91aa-1ae62193fa54 2
			`),
			WantMarkdown: new(`
            | ID                                   | ROLE                | INHERITED | SINCE |
            |--------------------------------------|---------------------|-----------|-------|
            | 16d6e8ba-f574-494f-8d5e-74f6cb2d8db0 | PROJECT_ROLE_OWNER  | false     | now   |
            | 40c0da4b-9eb9-4371-91aa-1ae62193fa54 | PROJECT_ROLE_EDITOR | true      | now   |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_DeleteMember(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceRemoveMemberResponse, string]{
		{
			Name:    "delete project member",
			CmdArgs: []string{"project", "member", "delete", testresources.ProjectMember1().Id, "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("RemoveMember", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceRemoveMemberRequest{
								Project:  testresources.Project1().Uuid,
								MemberId: testresources.ProjectMember1().Id,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceRemoveMemberResponse{}), nil)
						},
					},
				}}),
			WantDefault: new(fmt.Sprintf("✔ successfully removed member \"%s\"", testresources.ProjectMember1().Id)),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_ProjectCmd_UpdateMember(t *testing.T) {
	tests := []*e2e.Test[apiv1.ProjectServiceUpdateMemberResponse, *apiv1.ProjectMember]{
		{
			Name:    "update project member",
			CmdArgs: []string{"project", "member", "update", testresources.ProjectMember1().Id, "--project", testresources.Project1().Uuid, "--role", testresources.ProjectMember1().Role.String()},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Project: func(m *mock.Mock) {
							m.On("UpdateMember", mock.Anything, connect.NewRequest(&apiv1.ProjectServiceUpdateMemberRequest{
								Project:  testresources.Project1().Uuid,
								MemberId: testresources.ProjectMember1().Id,
								Role:     testresources.ProjectMember1().Role,
							})).Return(connect.NewResponse(&apiv1.ProjectServiceUpdateMemberResponse{
								ProjectMember: testresources.ProjectMember1(),
							}), nil)
						},
					},
				}}),
			WantObject: testresources.ProjectMember1(),
			WantTable: new(`
			ID                                    ROLE                INHERITED  SINCE  
            16d6e8ba-f574-494f-8d5e-74f6cb2d8db0  PROJECT_ROLE_OWNER  false      now
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
