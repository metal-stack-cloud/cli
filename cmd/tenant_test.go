package cmd

import (
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_TenantCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.Tenant]{
		{
			Name: "list tenants",
			Cmd: func(want []*apiv1.Tenant) []string {
				return []string{"admin", "tenant", "list"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Adminv1Mocks: &apitests.Adminv1MockFns{
					Tenant: func(m *mock.Mock) {
						m.On("List", mock.Anything, connect.NewRequest(&adminv1.TenantServiceListRequest{})).Return(&connect.Response[adminv1.TenantServiceListResponse]{
							Msg: &adminv1.TenantServiceListResponse{
								Tenants: []*apiv1.Tenant{
									{
										Login:         "loginOne",
										Name:          "nameOne",
										Email:         "testone@mail.com",
										OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
										Admitted:      false,
										CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Hour)),
									},
									{
										Login:         "loginTwo",
										Name:          "nameTwo",
										Email:         "testtwo@mail.com",
										OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
										Admitted:      true,
										CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Minute)),
									},
								},
							},
						}, nil)
					},
				},
			},
			Want: []*apiv1.Tenant{
				{
					Login:         "loginTwo",
					Name:          "nameTwo",
					Email:         "testtwo@mail.com",
					OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
					Admitted:      true,
					CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Minute)),
				},
				{
					Login:         "loginOne",
					Name:          "nameOne",
					Email:         "testone@mail.com",
					OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
					Admitted:      false,
					CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Hour)),
				},
			},
			WantTable: pointer.Pointer(`
ID         NAME      EMAIL              PROVIDER                 REGISTERED     ADMITTED   COUPONS
loginTwo   nameTwo   testtwo@mail.com   O_AUTH_PROVIDER_AZURE    1 minute ago   true       -
loginOne   nameOne   testone@mail.com   O_AUTH_PROVIDER_GITHUB   1 hour ago     false      -
`),
			WantWideTable: pointer.Pointer(`
ID         NAME      EMAIL              PROVIDER                 REGISTERED     ADMITTED   COUPONS
loginTwo   nameTwo   testtwo@mail.com   O_AUTH_PROVIDER_AZURE    1 minute ago   true       -
loginOne   nameOne   testone@mail.com   O_AUTH_PROVIDER_GITHUB   1 hour ago     false      -
`),
			Template: pointer.Pointer("{{ .login }} {{ .name }} {{ .email }} {{ .oauth_provider }} {{ if .admitted }}true{{ else }}false{{ end }}"),
			WantTemplate: pointer.Pointer(`
loginTwo nameTwo testtwo@mail.com 2 true
loginOne nameOne testone@mail.com 1 false
			`),
			WantMarkdown: pointer.Pointer(`
|    ID    |  NAME   |      EMAIL       |        PROVIDER        |  REGISTERED  | ADMITTED | COUPONS |
|----------|---------|------------------|------------------------|--------------|----------|---------|
| loginTwo | nameTwo | testtwo@mail.com | O_AUTH_PROVIDER_AZURE  | 1 minute ago | true     | -       |
| loginOne | nameOne | testone@mail.com | O_AUTH_PROVIDER_GITHUB | 1 hour ago   | false    | -       |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
func Test_TenantCmd_SingleResult(t *testing.T) {
	tests := []*Test[*apiv1.Tenant]{
		{
			Name: "admit",
			Cmd: func(want *apiv1.Tenant) []string {
				return []string{"admin", "tenant", "admit", "someid"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Adminv1Mocks: &apitests.Adminv1MockFns{
					Tenant: func(m *mock.Mock) {
						m.On("Admit", mock.Anything, connect.NewRequest(&adminv1.TenantServiceAdmitRequest{
							TenantId: "someid",
						})).Return(&connect.Response[adminv1.TenantServiceAdmitResponse]{
							Msg: &adminv1.TenantServiceAdmitResponse{
								Tenant: &apiv1.Tenant{
									Login:         "someid",
									Name:          "somename",
									Email:         "somename@mail.com",
									AvatarUrl:     "https://avatars.githubusercontent.com/u/52534857?v=8",
									OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
									Admitted:      true,
									PhoneNumber:   "1234556",
									EmailConsent:  true,
									CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Minute)),
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.Tenant{
				Login:         "someid",
				Name:          "somename",
				Email:         "somename@mail.com",
				AvatarUrl:     "https://avatars.githubusercontent.com/u/52534857?v=8",
				OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
				Admitted:      true,
				PhoneNumber:   "1234556",
				EmailConsent:  true,
				CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Minute)),
			},
		},
		{
			Name: "revoke",
			Cmd: func(want *apiv1.Tenant) []string {
				return []string{"admin", "tenant", "revoke", "someid"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Adminv1Mocks: &apitests.Adminv1MockFns{
					Tenant: func(m *mock.Mock) {
						m.On("Revoke", mock.Anything, connect.NewRequest(&adminv1.TenantServiceRevokeRequest{
							TenantId: "someid",
						})).Return(&connect.Response[adminv1.TenantServiceRevokeResponse]{
							Msg: &adminv1.TenantServiceRevokeResponse{
								Tenant: &apiv1.Tenant{
									Login:         "someid",
									Name:          "somename",
									Email:         "somename@mail.com",
									AvatarUrl:     "https://avatars.githubusercontent.com/u/52534857?v=8",
									OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
									PhoneNumber:   "1234556",
									EmailConsent:  true,
									CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Minute)),
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.Tenant{
				Login:         "someid",
				Name:          "somename",
				Email:         "somename@mail.com",
				AvatarUrl:     "https://avatars.githubusercontent.com/u/52534857?v=8",
				OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_AZURE,
				PhoneNumber:   "1234556",
				EmailConsent:  true,
				CreatedAt:     timestamppb.New(time.Now().Add(-1 * time.Minute)),
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
