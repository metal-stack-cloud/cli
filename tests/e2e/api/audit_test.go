package api_e2e

import (
	"strconv"
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/stretchr/testify/mock"
)

func Test_AuditCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.AuditServiceGetResponse, *apiv1.AuditTrace]{
		{
			Name:    "describe audit trace",
			CmdArgs: []string{"audit", "describe", testresources.AuditTrace1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						User: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.UserServiceGetRequest{})).
								Return(connect.NewResponse(&apiv1.UserServiceGetResponse{
									User: testresources.User1(),
								}), nil)
						},
						Audit: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.AuditServiceGetRequest{
								Login: testresources.AuditTrace1().Tenant,
								Uuid:  testresources.AuditTrace1().Uuid,
							})).
								Return(connect.NewResponse(&apiv1.AuditServiceGetResponse{
									Trace: testresources.AuditTrace1(),
								}), nil)
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
	tests := []*e2e.Test[apiv1.AuditServiceListResponse, []*apiv1.AuditTrace]{
		{
			Name:    "list",
			CmdArgs: []string{"audit", "list", "--tenant", testresources.AuditTrace1().Tenant},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Audit: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.AuditServiceListRequest{
								Login: testresources.AuditTrace1().Tenant,
							})).
								Return(connect.NewResponse(&apiv1.AuditServiceListResponse{
									Traces: []*apiv1.AuditTrace{
										testresources.AuditTrace1(),
									},
								}), nil)
						},
					},
				},
			}),
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
		{
			Name: "list with filters",
			CmdArgs: []string{"audit", "list",
				"--tenant", testresources.AuditTrace1().Tenant,
				"--request-id", testresources.AuditTrace1().Uuid,
				"--from", testresources.AuditTrace1().Timestamp.AsTime().Format("2006-01-02 15:04:05"),
				"--to", testresources.AuditTrace1().Timestamp.AsTime().Format("2006-01-02 15:04:05"),
				"--user", testresources.AuditTrace1().User,
				"--project", *testresources.AuditTrace1().Project,
				"--method", testresources.AuditTrace1().Method,
				"--source-ip", testresources.AuditTrace1().SourceIp,
				"--result-code", strconv.Itoa(int(*testresources.AuditTrace1().ResultCode)),
				"--limit", "100",
				"--phase", testresources.AuditTrace1().Phase.String(),
				"--body", *testresources.AuditTrace1().Body,
				"--prettify-body",
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Audit: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.AuditServiceListRequest{
								Login:      testresources.AuditTrace1().Tenant,
								Uuid:       &testresources.AuditTrace1().Uuid,
								From:       testresources.AuditTrace1().Timestamp,
								To:         testresources.AuditTrace1().Timestamp,
								User:       &testresources.AuditTrace1().User,
								Project:    testresources.AuditTrace1().Project,
								Method:     &testresources.AuditTrace1().Method,
								ResultCode: testresources.AuditTrace1().ResultCode,
								SourceIp:   &testresources.AuditTrace1().SourceIp,
								Body:       testresources.AuditTrace1().Body,
								Limit:      new(int32(100)),
								Phase:      &testresources.AuditTrace1().Phase,
							})).
								Return(connect.NewResponse(&apiv1.AuditServiceListResponse{
									Traces: []*apiv1.AuditTrace{
										testresources.AuditTrace1(),
									},
								}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            TIME                 REQUEST - ID                          USER         PROJECT                               METHOD     PHASE                CODE  
            2000-01-01 00:00:00  c40ad996-e1fd-4511-a7bf-418219cb8d91  metal-stack  c40ad996-e1fd-4511-a7bf-418219cb8d95  /apiv1/ip  AUDIT_PHASE_REQUEST  OK
				`),
			WantWideTable: new(`
            TIME                 REQUEST - ID                          USER         PROJECT                               METHOD     PHASE                SOURCE - IP  CODE  BODY          
            2000-01-01 00:00:00  c40ad996-e1fd-4511-a7bf-418219cb8d91  metal-stack  c40ad996-e1fd-4511-a7bf-418219cb8d95  /apiv1/ip  AUDIT_PHASE_REQUEST  192.168.2.1  OK    {             
                                                                                                                                                                                 "a": "b"  
                                                                                                                                                                             }
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
