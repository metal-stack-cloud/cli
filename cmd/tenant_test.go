package cmd

import (
	"testing"
)

func Test_TenantCmd_MultiResult(t *testing.T) {
// 	tests := []*Test[[]*apiv1.Tenant]{
// 		{
// 			Name: "list",
// 			Cmd: func(want []*apiv1.Tenant) []string {
// 				return []string{"admin", "tenant", "list"}
// 			},
// 			AdminMocks: &apitests.AdminMockFns{
// 				Tenant: func(m *mock.Mock) {
// 					m.On("List", mock.Anything, connect.NewRequest(&adminv1.TenantServiceListRequest{})).Return(&connect.Response[adminv1.TenantServiceListResponse]{
// 						Msg: &adminv1.TenantServiceListResponse{
// 							Tenants: []*apiv1.Tenant{
// 								{
// 									Login: "loginOne",
// 									Name: "nameOne",
// 									Email: "testone@mail.com",
// 									OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
// 									Admitted: false,
// 								},
// 								{
// 									Login: "loginTwo",
// 									Name: "nameTwo",
// 									Email: "testtwo@mail.com",
// 									OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
// 									Admitted: true,
// 								},
// 							},
// 						},
// 					}, nil)
// 				},
// 			},
// 			Want: []*apiv1.Tenant{
// 				{
// 					Login: "loginOne",
// 					Name: "nameOne",
// 					Email: "testone@mail.com",
// 					OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
// 					Admitted: false,
// 				},
// 				{
// 					Login: "loginTwo",
// 					Name: "nameTwo",
// 					Email: "testtwo@mail.com",
// 					OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
// 					Admitted: true,
// 				},
// },
// 			WantTable: pointer.Pointer(`
// ID         NAME      EMAIL             PROVIDER                 REGISTERED         ADMITTED
// loginOne   nameOne   testone@mail.de   O_AUTH_PROVIDER_GITHUB   a long while ago   false
// loginTwo   nameTwo   testtwo@mail.de   O_AUTH_PROVIDER_AZURE    a long while ago   true
// `),
// 			WantWideTable: pointer.Pointer(`
// ID         NAME      EMAIL             PROVIDER                 REGISTERED         ADMITTED
// loginOne   nameOne   testone@mail.de   O_AUTH_PROVIDER_GITHUB   a long while ago   false
// loginTwo   nameTwo   testtwo@mail.de   O_AUTH_PROVIDER_AZURE    a long while ago   true
// `),
// 			Template: pointer.Pointer("{{ .id }} {{ .name }} {{ .email }} {{ .provider }} {{ .registered }} {{ .admitted }}"),
// 			WantTemplate: pointer.Pointer(`
// loginOne nameOne testone@mail.de O_AUTH_PROVIDER_GITHUB a long while ago false
// loginTwo nameTwo testtwo@mail.de O_AUTH_PROVIDER_AZURE a long while ago true
// 			`),
// 			WantMarkdown: pointer.Pointer(`
// |    ID    |   NAME  | EMAIL           |  PROVIDER              |  REGISTERED      |  ADMITTED |
// |----------|---------|-----------------|------------------------|------------------|-----------|
// | loginOne | nameOne | testone@mail.de | O_AUTH_PROVIDER_GITHUB | a long while ago | false     |
// | loginTwo | nameTwo | testtwo@mail.de | O_AUTH_PROVIDER_AZURE  | a long while ago | true      |
// `),
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt.TestCmd(t)
// 	}
}
