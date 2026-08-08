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

func Test_AdminTenantCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.TenantServiceListResponse, []*apiv1.Tenant]{
		{
			Name:    "list tenants",
			CmdArgs: []string{"admin", "tenant", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&adminv1.TenantServiceListRequest{})).Return(connect.NewResponse(&adminv1.TenantServiceListResponse{
								Tenants: []*apiv1.Tenant{
									testresources.Tenant1(), testresources.Tenant2(),
								},
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
			WantWideTable: new(`
            ID             NAME         EMAIL                 REGISTERED  PROVIDER                ADMITTED  TERMS AND CONDITIONS  
            metal-stack    metal-stack  metal-stack@mail.com  now         O_AUTH_PROVIDER_GITHUB  true      true                  
            x-cellent.com  x-cellent    x-cellent@mail.com    now         O_AUTH_PROVIDER_GITHUB  false     false
`),
			Template: new("{{ .login }} {{ .name }} {{ .email }} {{ .oauth_provider }} {{ if .admitted }}true{{ else }}false{{ end }} {{ if .terms_and_conditions.accepted }}true{{ else }}false{{ end }}"),
			WantTemplate: new(`
metal-stack metal-stack metal-stack@mail.com 1 true true
x-cellent.com x-cellent x-cellent@mail.com 1 false false
			`),
			WantMarkdown: new(`
            | ID            | NAME        | EMAIL                | REGISTERED | PROVIDER               | ADMITTED | TERMS AND CONDITIONS |
            |---------------|-------------|----------------------|------------|------------------------|----------|----------------------|
            | metal-stack   | metal-stack | metal-stack@mail.com | now        | O_AUTH_PROVIDER_GITHUB | true     | true                 |
            | x-cellent.com | x-cellent   | x-cellent@mail.com   | now        | O_AUTH_PROVIDER_GITHUB | false    | false                |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTenantCmd_Admit(t *testing.T) {
	tests := []*e2e.Test[adminv1.TenantServiceAdmitResponse, *apiv1.Tenant]{
		{
			Name:    "admit",
			CmdArgs: []string{"admin", "tenant", "admit", testresources.Tenant2().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Admit", mock.Anything, connect.NewRequest(&adminv1.TenantServiceAdmitRequest{
								TenantId: testresources.Tenant2().Login,
							})).Return(connect.NewResponse(&adminv1.TenantServiceAdmitResponse{
								Tenant: testresources.Tenant2(),
							}), nil)
						},
					},
				}}),

			WantObject: testresources.Tenant2(),
			WantDefault: new(`
createdAt: "2000-01-01T00:00:00Z"
createdBy: metal-stack
description: test Tenant 2
email: x-cellent@mail.com
login: x-cellent.com
name: x-cellent
oauthProvider: O_AUTH_PROVIDER_GITHUB
onboarded: true
phoneNumber: "0123456787"
termsAndConditions: {}
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTenantCmd_Revoke(t *testing.T) {
	tests := []*e2e.Test[adminv1.TenantServiceRevokeResponse, *apiv1.Tenant]{
		{
			Name:    "revoke",
			CmdArgs: []string{"admin", "tenant", "revoke", testresources.Tenant2().Login},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("Revoke", mock.Anything, connect.NewRequest(&adminv1.TenantServiceRevokeRequest{
								TenantId: testresources.Tenant2().Login,
							})).Return(connect.NewResponse(&adminv1.TenantServiceRevokeResponse{
								Tenant: testresources.Tenant2(),
							}), nil)
						},
					},
				}}),

			WantObject: testresources.Tenant2(),
			WantDefault: new(`
createdAt: "2000-01-01T00:00:00Z"
createdBy: metal-stack
description: test Tenant 2
email: x-cellent@mail.com
login: x-cellent.com
name: x-cellent
oauthProvider: O_AUTH_PROVIDER_GITHUB
onboarded: true
phoneNumber: "0123456787"
termsAndConditions: {}
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminTenantCmd_AddMember(t *testing.T) {
	tests := []*e2e.Test[adminv1.TenantServiceAddMemberResponse, string]{
		{
			Name: "add Member",
			CmdArgs: []string{"admin", "tenant", "add-member",
				"--tenant-id", testresources.Tenant1().Login,
				"--member-id", testresources.TenantMember1().Id,
				"--role", testresources.TenantMember1().Role.String(),
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("AddMember", mock.Anything, connect.NewRequest(&adminv1.TenantServiceAddMemberRequest{
								TenantId: testresources.Tenant1().Login,
								MemberId: testresources.TenantMember1().Id,
								Role:     testresources.TenantMember1().Role,
							})).Return(connect.NewResponse(&adminv1.TenantServiceAddMemberResponse{}), nil)
						},
					},
				},
			}),
			WantDefault: new(``),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
