package testresources

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Tenant1 = func() *apiv1.Tenant {
		return &apiv1.Tenant{
			Login:         "metal-stack",
			Name:          "metal-stack",
			Email:         "metal-stack@mail.com",
			Description:   "test Tenant 1",
			OauthProvider: apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
			Admitted:      true,
			PhoneNumber:   "0123456789",
			EmailConsent:  true,
			Onboarded:     true,
			CreatedAt:     timestamppb.New(e2e.TimeBubbleStartTime()),
			PaymentDetails: &apiv1.PaymentDetails{
				CustomerId:      "cus-some1234",
				PaymentMethodId: new("millionsofdollars"),
				SubscriptionId:  "sub-123456",
			},
			TermsAndConditions: &apiv1.TermsAndConditions{
				Accepted: true,
			},
		}
	}
	Tenant2 = func() *apiv1.Tenant {
		return &apiv1.Tenant{
			Login:          "x-cellent.com",
			Name:           "x-cellent",
			Email:          "x-cellent@mail.com",
			Description:    "test Tenant 2",
			OauthProvider:  apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
			Admitted:       false,
			PhoneNumber:    "0123456787",
			EmailConsent:   false,
			Onboarded:      true,
			CreatedBy:      Tenant1().Name,
			CreatedAt:      timestamppb.New(e2e.TimeBubbleStartTime()),
			PaymentDetails: nil,
			TermsAndConditions: &apiv1.TermsAndConditions{
				Accepted: false,
			},
		}
	}
	TenantMember1 = func() *apiv1.TenantMember {
		return &apiv1.TenantMember{
			Id:         User1().Login,
			Role:       apiv1.TenantRole_TENANT_ROLE_OWNER,
			ProjectIds: []string{Project1().Uuid, Project2().Uuid},
			CreatedAt:  timestamppb.New(e2e.TimeBubbleStartTime()),
		}
	}
	TenantMember2 = func() *apiv1.TenantMember {
		return &apiv1.TenantMember{
			Id:         User2().Login,
			Role:       apiv1.TenantRole_TENANT_ROLE_OWNER,
			ProjectIds: []string{Project2().Uuid},
			CreatedAt:  timestamppb.New(e2e.TimeBubbleStartTime()),
		}
	}
)
