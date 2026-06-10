package testresources

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	AuditTrace1 = func() *apiv1.AuditTrace {
		return &apiv1.AuditTrace{
			Uuid:       "c40ad996-e1fd-4511-a7bf-418219cb8d91",
			Timestamp:  timestamppb.New(e2e.TimeBubbleStartTime()),
			User:       User1().Login,
			Tenant:     Tenant1().Login,
			Project:    &Project1().Uuid,
			Method:     "/apiv1/ip",
			Body:       new(`{"a": "b"}`),
			SourceIp:   "192.168.2.1",
			ResultCode: new(int32(codes.OK)),
			Phase:      apiv1.AuditPhase_AUDIT_PHASE_REQUEST,
		}
	}
	AuditTrace2 = func() *apiv1.AuditTrace {
		return &apiv1.AuditTrace{
			Uuid:       "b5817ef7-980a-41ef-9ed3-741a143870b0",
			Timestamp:  timestamppb.New(e2e.TimeBubbleStartTime()),
			User:       User2().Login,
			Tenant:     Tenant2().Login,
			Project:    &Project2().Uuid,
			Method:     "/apiv1/cluster",
			Body:       new(`{"c": "d"}`),
			SourceIp:   "192.168.2.3",
			ResultCode: new(int32(codes.NotFound)),
			Phase:      apiv1.AuditPhase_AUDIT_PHASE_RESPONSE,
		}
	}
)
