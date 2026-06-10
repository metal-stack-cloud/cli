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

var (
	health1 = func() *apiv1.Health {
		return &apiv1.Health{
			Services: []*apiv1.HealthStatus{
				{
					Name:    apiv1.Service_SERVICE_CLUSTER,
					Status:  apiv1.ServiceStatus_SERVICE_STATUS_HEALTHY,
					Message: "i am healthy",
				},
				{
					Name:    apiv1.Service_SERVICE_MACHINES,
					Status:  apiv1.ServiceStatus_SERVICE_STATUS_UNHEALTHY,
					Message: "i am not healthy",
				},
				{
					Name:    apiv1.Service_SERVICE_VOLUME,
					Status:  apiv1.ServiceStatus_SERVICE_STATUS_DEGRADED,
					Message: "i am",
				},
			},
		}
	}
)

func Test_HealthCmd(t *testing.T) {
	tests := []*e2e.Test[apiv1.HealthServiceGetResponse, *apiv1.Health]{
		{
			Name:    "health",
			CmdArgs: []string{"health"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Health: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.HealthServiceGetRequest{})).Return(connect.NewResponse(&apiv1.HealthServiceGetResponse{
								Health: health1(),
							}), nil)
						},
					},
				},
			}),
			WantDefault: new(`
NAME      MESSAGE           
✔  Clusters  i am healthy      
✗  Machines  i am not healthy  
✗  Volumes   i am
			`),
			WantMarkdown: new(`
            |   | NAME     | MESSAGE          |
            |---|----------|------------------|
            | ✔ | Clusters | i am healthy     |
            | ✗ | Machines | i am not healthy |
            | ✗ | Volumes  | i am             |
			`),
			WantObject:      health1(),
			WantProtoObject: health1(),
			Template:        new("{{ range $s := .services }}{{ $s.message }} {{ end }}"),
			WantTemplate:    new(`i am healthy i am not healthy i am`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
