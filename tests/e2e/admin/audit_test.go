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

func Test_AuditCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[adminv1.AuditServiceGetResponse, *apiv1.AuditTrace]{
		{
			Name:    "describe audit trace",
			CmdArgs: []string{"admin", "audit", "describe", testresources.AuditTrace1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Audit: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&adminv1.AuditServiceGetRequest{
								Uuid: testresources.AuditTrace1().Uuid,
							})).
								Return(&connect.Response[adminv1.AuditServiceGetResponse]{
									Msg: &adminv1.AuditServiceGetResponse{
										Trace: testresources.AuditTrace1(),
									},
								}, nil)
						},
					},
				},
			}),
			WantObject: testresources.AuditTrace1(),
			WantTable: new(`
         	TIME                 REQUEST - ID                          USER         PROJECT                               METHOD     PHASE                CODE  
            2000-01-01 00:00:00  c40ad996-e1fd-4511-a7bf-418219cb8d91  metal-stack  c40ad996-e1fd-4511-a7bf-418219cb8d95  /apiv1/ip  AUDIT_PHASE_REQUEST  OK
				`),
			WantWideTable: new(`
            TIME                 REQUEST - ID                          USER         PROJECT                               METHOD     PHASE                SOURCE - IP  CODE  BODY        
            2000-01-01 00:00:00  c40ad996-e1fd-4511-a7bf-418219cb8d91  metal-stack  c40ad996-e1fd-4511-a7bf-418219cb8d95  /apiv1/ip  AUDIT_PHASE_REQUEST  192.168.2.1  OK    {"a": "b"}
				`),
			Template: new(`{{ date "02/01/2006" .timestamp }} {{ .uuid }}`),
			WantTemplate: new(`
01/01/2000 c40ad996-e1fd-4511-a7bf-418219cb8d91
				`),
			WantMarkdown: new(`
            | TIME                | REQUEST - ID                         | USER        | PROJECT                              | METHOD    | PHASE               | CODE |
            |---------------------|--------------------------------------|-------------|--------------------------------------|-----------|---------------------|------|
            | 2000-01-01 00:00:00 | c40ad996-e1fd-4511-a7bf-418219cb8d91 | metal-stack | c40ad996-e1fd-4511-a7bf-418219cb8d95 | /apiv1/ip | AUDIT_PHASE_REQUEST | OK   |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AuditCmd_List(t *testing.T) {
	tests := []*e2e.Test[adminv1.AuditServiceListResponse, []*apiv1.AuditTrace]{
		{
			Name:    "list",
			CmdArgs: []string{"admin", "audit", "list", "--tenant", testresources.AuditTrace1().Tenant},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Adminv1Mocks: &apitests.Adminv1MockFns{
						Audit: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&adminv1.AuditServiceListRequest{})).
								Return(&connect.Response[adminv1.AuditServiceListResponse]{
									Msg: &adminv1.AuditServiceListResponse{
										Traces: []*apiv1.AuditTrace{
											testresources.AuditTrace1(),
											testresources.AuditTrace2(),
										},
									},
								}, nil)
						},
					},
				},
			}),
			WantTable: new(`
            TIME                 REQUEST - ID                          USER           PROJECT                               METHOD          PHASE                 CODE      
            2000-01-01 00:00:00  c40ad996-e1fd-4511-a7bf-418219cb8d91  metal-stack    c40ad996-e1fd-4511-a7bf-418219cb8d95  /apiv1/ip       AUDIT_PHASE_REQUEST   OK        
            2000-01-01 00:00:00  b5817ef7-980a-41ef-9ed3-741a143870b0  x-cellent.com  c40ad996-e1fd-4511-a7bf-418219cb8d67  /apiv1/cluster  AUDIT_PHASE_RESPONSE  NotFound
				`),
			WantWideTable: new(`
            TIME                 REQUEST - ID                          USER           PROJECT                               METHOD          PHASE                 SOURCE - IP  CODE      BODY        
            2000-01-01 00:00:00  c40ad996-e1fd-4511-a7bf-418219cb8d91  metal-stack    c40ad996-e1fd-4511-a7bf-418219cb8d95  /apiv1/ip       AUDIT_PHASE_REQUEST   192.168.2.1  OK        {"a": "b"}  
            2000-01-01 00:00:00  b5817ef7-980a-41ef-9ed3-741a143870b0  x-cellent.com  c40ad996-e1fd-4511-a7bf-418219cb8d67  /apiv1/cluster  AUDIT_PHASE_RESPONSE  192.168.2.3  NotFound  {"c": "d"}
				`),
			Template: new(`{{ date "02/01/2006" .timestamp }} {{ .uuid }}`),
			WantTemplate: new(`
01/01/2000 c40ad996-e1fd-4511-a7bf-418219cb8d91
01/01/2000 b5817ef7-980a-41ef-9ed3-741a143870b0
				`),
			WantMarkdown: new(`
            | TIME                | REQUEST - ID                         | USER          | PROJECT                              | METHOD         | PHASE                | CODE     |
            |---------------------|--------------------------------------|---------------|--------------------------------------|----------------|----------------------|----------|
            | 2000-01-01 00:00:00 | c40ad996-e1fd-4511-a7bf-418219cb8d91 | metal-stack   | c40ad996-e1fd-4511-a7bf-418219cb8d95 | /apiv1/ip      | AUDIT_PHASE_REQUEST  | OK       |
            | 2000-01-01 00:00:00 | b5817ef7-980a-41ef-9ed3-741a143870b0 | x-cellent.com | c40ad996-e1fd-4511-a7bf-418219cb8d67 | /apiv1/cluster | AUDIT_PHASE_RESPONSE | NotFound |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
