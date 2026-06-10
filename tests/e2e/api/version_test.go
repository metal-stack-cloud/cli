package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/stretchr/testify/mock"
)

func Test_VersionCmd(t *testing.T) {
	tests := []*e2e.Test[apiv1.VersionServiceGetResponse, *apiv1.Version]{
		{
			Name:    "get Version",
			CmdArgs: []string{"version"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Version: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VersionServiceGetRequest{})).Return(&connect.Response[apiv1.VersionServiceGetResponse]{
								Msg: &apiv1.VersionServiceGetResponse{
									Version: testresources.Version(),
								},
							}, nil)
						},
					},
				},
			}),
			Template:     new("{{ .Server.build_date }} {{ .Server.git_sha1 }} {{ .Server.revision }} {{ .Server.version }}"),
			WantTemplate: new(`2026-03-21T15:35:07+00:00 477edc0b tags/v0.1.8-0-g476edc0 v0.1.8`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
