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

func Test_TenantCmd_Get(t *testing.T) {
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

// doesnt work solo
func Test_TenantCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.TenantServiceListResponse, []*apiv1.Tenant]{
		{
			Name:    "list tenants",
			CmdArgs: []string{"tenant", "list"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Tenant: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.TenantServiceListRequest{
								Id: &testresources.Tenant1().Login,
							})).Return(connect.NewResponse(&apiv1.TenantServiceListResponse{
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
