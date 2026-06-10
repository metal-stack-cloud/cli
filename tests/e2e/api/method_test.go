package api_e2e

import (
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/stretchr/testify/mock"
)

func Test_MethodCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.MethodServiceListResponse, []string]{
		{
			Name:    "list api methods",
			CmdArgs: []string{"api-methods"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Method: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.MethodServiceListRequest{})).Return(connect.NewResponse(&apiv1.MethodServiceListResponse{
								Methods: []string{},
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(``),
		},
		{
			Name:    "list api methods",
			CmdArgs: []string{"api-methods", "--scoped", "true"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Method: func(m *mock.Mock) {
							m.On("TokenScopedList", mock.Anything, connect.NewRequest(&apiv1.MethodServiceTokenScopedListRequest{})).Return(connect.NewResponse(&apiv1.MethodServiceTokenScopedListResponse{}), nil)
						},
					},
				},
			}),
			WantDefault: new(`{}`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
