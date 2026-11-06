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

func Test_AssetCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.AssetServiceListResponse, []*apiv1.Asset]{
		{
			Name:    "list assets",
			CmdArgs: []string{"asset"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Asset: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.AssetServiceListRequest{})).Return(connect.NewResponse(&apiv1.AssetServiceListResponse{
								Assets: []*apiv1.Asset{testresources.Asset1()},
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(`
REGION  PARTITION  MACHINE TYPES  
fra     fra-1      c1-medium-x86  
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
