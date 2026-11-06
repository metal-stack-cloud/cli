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

func Test_UserCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.UserServiceGetResponse, *apiv1.User]{
		{
			Name:    "user describe",
			CmdArgs: []string{"user", "describe"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						User: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.UserServiceGetRequest{})).Return(connect.NewResponse(&apiv1.UserServiceGetResponse{
								User: testresources.User1(),
							}), nil)
						},
					},
				}}),

			WantObject:      testresources.User1(),
			WantProtoObject: testresources.User1(),
			WantTable: new(`
            LOGIN        NAME   EMAIL               
            metal-stack  Larry  evilLarry@mail.com
`),
			WantWideTable: new(`
            LOGIN        NAME   EMAIL               TENANTS                 PROJECTS                                   
            metal-stack  Larry  evilLarry@mail.com  metal-stack, x-cellent  Some Initiative, Some Initiative Number 2
`),
			WantMarkdown: new(`
            | LOGIN       | NAME  | EMAIL              |
            |-------------|-------|--------------------|
            | metal-stack | Larry | evilLarry@mail.com |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
