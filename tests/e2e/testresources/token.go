package testresources

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Token1 = func() *apiv1.Token {
		return &apiv1.Token{
			Uuid:        "123e4567-e89b-12d3-a456-426614174000",
			UserId:      "user_987654321",
			Description: "Production API read-only token for CI/CD pipeline",
			Permissions: []*apiv1.MethodPermission{
				{
					Subject: "project-01",
					Methods: []string{"api/ip", "api/health"},
				},
			},
			Expires:   timestamppb.New(e2e.TimeBubbleStartTime().Add(24 * time.Hour)),
			IssuedAt:  timestamppb.New(e2e.TimeBubbleStartTime()),
			TokenType: apiv1.TokenType_TOKEN_TYPE_API,
			ProjectRoles: map[string]apiv1.ProjectRole{
				"project-01": apiv1.ProjectRole_PROJECT_ROLE_EDITOR,
			},
			TenantRoles: map[string]apiv1.TenantRole{
				"tenant-01": apiv1.TenantRole_TENANT_ROLE_OWNER,
			},
			AdminRole: new(apiv1.AdminRole_ADMIN_ROLE_EDITOR),
		}
	}
	Token2 = func() *apiv1.Token {
		return &apiv1.Token{
			Uuid:        "ccd042af-7e3a-4458-bc9f-933a26fd783e",
			UserId:      "user_Guest",
			Description: "GuestToken",
			Permissions: []*apiv1.MethodPermission{
				{
					Subject: "tenant-01",
					Methods: []string{"api/ip", "api/health"},
				},
			},
			Expires:   timestamppb.New(e2e.TimeBubbleStartTime().Add(24 * time.Hour)),
			IssuedAt:  timestamppb.New(e2e.TimeBubbleStartTime()),
			TokenType: apiv1.TokenType_TOKEN_TYPE_CONSOLE,
			TenantRoles: map[string]apiv1.TenantRole{
				"tenant-01": apiv1.TenantRole_TENANT_ROLE_GUEST,
			},
		}
	}
)
