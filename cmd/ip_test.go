package cmd

import (
	"testing"

	"github.com/bufbuild/connect-go"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/mock"
)

func Test_IPCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.IP]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.IP) []string {
				return []string{"ip", "list", "--project", "a"}
			},
			APIMocks: &apitests.APIMockFns{
				IP: func(m *mock.Mock) {
					m.On("List", mock.Anything, connect.NewRequest(&apiv1.IPServiceListRequest{
						Project: "a",
					})).Return(&connect.Response[apiv1.IPServiceListResponse]{
						Msg: &apiv1.IPServiceListResponse{
							Ips: []*apiv1.IP{
								{
									Uuid:        "uuid",
									Ip:          "1.1.1.1",
									Name:        "a",
									Description: "a description",
									Network:     "a-network",
									Project:     "a",
									Type:        "ephemeral",
									Tags:        []string{},
								},
							},
						},
					}, nil)
				},
			},
			Want: []*apiv1.IP{
				{
					Ip: "1.1.1.1",
				},
			},
			WantTable: pointer.Pointer(`

`),
			WantWideTable: pointer.Pointer(`
`),
			// Template: pointer.Pointer("{{ .id }} {{ .name }}"),
			// 			WantTemplate: pointer.Pointer(`
			// 1 firewall-1
			// 2 firewall-2
			// `),
			WantMarkdown: pointer.Pointer(`
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
