package cmd

import (
	"net/http"
	"strconv"
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	code1       = int32(http.StatusOK)
	auditTrace1 = &apiv1.AuditTrace{
		Uuid:       "c40ad996-e1fd-4511-a7bf-418219cb8d91",
		Timestamp:  timestamppb.New(testTime),
		User:       "a-user",
		Tenant:     "a-tenant",
		Project:    pointer.Pointer("project-a"),
		Method:     "/apiv1/ip",
		Body:       pointer.Pointer(`{"a": "b"}`),
		SourceIp:   "192.168.2.1",
		ResultCode: &code1,
		Phase:      apiv1.AuditPhase_AUDIT_PHASE_REQUEST,
	}
	projectB    = "project-b"
	body2       = `{"c": "d"}`
	code2       = int32(http.StatusForbidden)
	auditTrace2 = &apiv1.AuditTrace{
		Uuid:       "b5817ef7-980a-41ef-9ed3-741a143870b0",
		Timestamp:  timestamppb.New(testTime),
		User:       "b-user",
		Tenant:     "b-tenant",
		Project:    &projectB,
		Method:     "/apiv1/cluster",
		Body:       &body2,
		SourceIp:   "192.168.2.3",
		ResultCode: &code2,
		Phase:      apiv1.AuditPhase_AUDIT_PHASE_RESPONSE,
	}
)

func Test_AuditCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.AuditTrace]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.AuditTrace) []string {
				return []string{"audit", "list", "--tenant", "a-tenant"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Audit: func(m *mock.Mock) {
						m.On("List", mock.Anything, connect.NewRequest(&apiv1.AuditServiceListRequest{
							Login: "a-tenant",
						})).
							Return(&connect.Response[apiv1.AuditServiceListResponse]{
								Msg: &apiv1.AuditServiceListResponse{
									Traces: []*apiv1.AuditTrace{
										auditTrace2,
										auditTrace1,
									},
								},
							}, nil)
					},
				},
			},
			Want: []*apiv1.AuditTrace{
				auditTrace2,
				auditTrace1,
			},
			WantTable: pointer.Pointer(`
TIME                  REQUEST-ID                             USER     PROJECT     METHOD           PHASE                  CODE
2022-05-19 01:02:03   b5817ef7-980a-41ef-9ed3-741a143870b0   b-user   project-b   /apiv1/cluster   AUDIT_PHASE_RESPONSE   403
2022-05-19 01:02:03   c40ad996-e1fd-4511-a7bf-418219cb8d91   a-user   project-a   /apiv1/ip        AUDIT_PHASE_REQUEST    200
				`),
			WantWideTable: pointer.Pointer(`
TIME                  REQUEST-ID                             USER     PROJECT     METHOD           PHASE                  SOURCE-IP     CODE   BODY
2022-05-19 01:02:03   b5817ef7-980a-41ef-9ed3-741a143870b0   b-user   project-b   /apiv1/cluster   AUDIT_PHASE_RESPONSE   192.168.2.3   403    {"c": "d"}
2022-05-19 01:02:03   c40ad996-e1fd-4511-a7bf-418219cb8d91   a-user   project-a   /apiv1/ip        AUDIT_PHASE_REQUEST    192.168.2.1   200    {"a": "b"}
				`),
			Template: pointer.Pointer(`{{ date "02/01/2006" .timestamp }} {{ .uuid }}`),
			WantTemplate: pointer.Pointer(`
19/05/2022 b5817ef7-980a-41ef-9ed3-741a143870b0
19/05/2022 c40ad996-e1fd-4511-a7bf-418219cb8d91
				`),
			WantMarkdown: pointer.Pointer(`
|        TIME         |              REQUEST-ID              |  USER  |  PROJECT  |     METHOD     |        PHASE         | CODE |
|---------------------|--------------------------------------|--------|-----------|----------------|----------------------|------|
| 2022-05-19 01:02:03 | b5817ef7-980a-41ef-9ed3-741a143870b0 | b-user | project-b | /apiv1/cluster | AUDIT_PHASE_RESPONSE |  403 |
| 2022-05-19 01:02:03 | c40ad996-e1fd-4511-a7bf-418219cb8d91 | a-user | project-a | /apiv1/ip      | AUDIT_PHASE_REQUEST  |  200 |
			`),
		},
		{
			Name: "list with filters",
			Cmd: func(want []*apiv1.AuditTrace) []string {
				project := *want[0].Project
				code := *want[0].ResultCode
				body := *want[0].Body
				args := []string{"audit", "list",
					"--tenant", "a-tenant",
					"--request-id", want[0].Uuid,
					"--from", want[0].Timestamp.AsTime().Format("2006-01-02 15:04:05"),
					"--to", want[0].Timestamp.AsTime().Format("2006-01-02 15:04:05"),
					"--user", want[0].User,
					"--project", project,
					"--method", want[0].Method,
					"--source-ip", want[0].SourceIp,
					"--result-code", strconv.Itoa(int(code)),
					"--error", "an-error",
					"--limit", "100",
					"--phase", want[0].Phase.String(),
					"--body", body,
					"--prettify-body",
				}
				AssertExhaustiveArgs(t, args, "sort-by")
				return args
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Audit: func(m *mock.Mock) {
						limit := int32(100)
						ts := auditTrace1.Timestamp
						ts.Nanos = 0 // nano sec prec gets lost from command line

						m.On("List", mock.Anything, connect.NewRequest(&apiv1.AuditServiceListRequest{
							Login:      "a-tenant",
							Uuid:       &auditTrace1.Uuid,
							From:       ts,
							To:         ts,
							User:       &auditTrace1.User,
							Project:    auditTrace1.Project,
							Method:     &auditTrace1.Method,
							ResultCode: auditTrace1.ResultCode,
							SourceIp:   &auditTrace1.SourceIp,
							Body:       auditTrace1.Body,
							Limit:      &limit,
							Phase:      &auditTrace1.Phase,
						})).
							Return(&connect.Response[apiv1.AuditServiceListResponse]{
								Msg: &apiv1.AuditServiceListResponse{
									Traces: []*apiv1.AuditTrace{
										auditTrace1,
									},
								},
							}, nil)
					},
				},
			},
			Want: []*apiv1.AuditTrace{
				auditTrace1,
			},
			WantTable: pointer.Pointer(`
TIME                  REQUEST-ID                             USER     PROJECT     METHOD      PHASE                 CODE
2022-05-19 01:02:03   c40ad996-e1fd-4511-a7bf-418219cb8d91   a-user   project-a   /apiv1/ip   AUDIT_PHASE_REQUEST   200
			`),
			WantWideTable: pointer.Pointer(`
TIME                  REQUEST-ID                             USER     PROJECT     METHOD      PHASE                 SOURCE-IP     CODE   BODY
2022-05-19 01:02:03   c40ad996-e1fd-4511-a7bf-418219cb8d91   a-user   project-a   /apiv1/ip   AUDIT_PHASE_REQUEST   192.168.2.1   200    {
                                                                                                                                             "a": "b"
                                                                                                                                         }
			`),
			Template: pointer.Pointer(`{{ date "02/01/2006" .timestamp }} {{ .uuid }}`),
			WantTemplate: pointer.Pointer(`
19/05/2022 c40ad996-e1fd-4511-a7bf-418219cb8d91
						`),
			WantMarkdown: pointer.Pointer(`
|        TIME         |              REQUEST-ID              |  USER  |  PROJECT  |  METHOD   |        PHASE        | CODE |
|---------------------|--------------------------------------|--------|-----------|-----------|---------------------|------|
| 2022-05-19 01:02:03 | c40ad996-e1fd-4511-a7bf-418219cb8d91 | a-user | project-a | /apiv1/ip | AUDIT_PHASE_REQUEST |  200 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
