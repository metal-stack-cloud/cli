package testresources

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	Project1 = func() *apiv1.Project {
		return &apiv1.Project{
			Uuid:             "c40ad996-e1fd-4511-a7bf-418219cb8d95",
			Name:             "Some Initiative",
			Description:      "Internal research and development for something.",
			Tenant:           Tenant1().Name,
			IsDefaultProject: false,
			CreatedAt:        timestamppb.New(e2e.TimeBubbleStartTime()),
			AvatarUrl:        new("https://cdn.example.com/avatars/me.png"),
		}
	}
	Project2 = func() *apiv1.Project {
		return &apiv1.Project{
			Uuid:             "c40ad996-e1fd-4511-a7bf-418219cb8d67",
			Name:             "Some Initiative Number 2",
			Description:      "Internal research and development for something even more important.",
			Tenant:           Tenant2().Name,
			IsDefaultProject: true,
			CreatedAt:        timestamppb.New(e2e.TimeBubbleStartTime()),
			AvatarUrl:        new("https://cdn.example.com/avatars/face.png"),
		}
	}
	ProjectInvite1 = func() *apiv1.ProjectInvite {
		return &apiv1.ProjectInvite{
			Secret:      "secret",
			Project:     "0d81bca7-73f6-4da3-8397-4a8c52a0c583",
			Role:        apiv1.ProjectRole_PROJECT_ROLE_EDITOR,
			Joined:      false,
			ProjectName: Project1().Name,
			TenantName:  Project1().Tenant,
			ExpiresAt:   timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
		}
	}
	ProjectInvite2 = func() *apiv1.ProjectInvite {
		return &apiv1.ProjectInvite{
			Secret:      "secret",
			Project:     "f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c",
			Role:        apiv1.ProjectRole_PROJECT_ROLE_EDITOR,
			Joined:      false,
			ProjectName: Project1().Name,
			TenantName:  Project1().Tenant,
			ExpiresAt:   timestamppb.New(e2e.TimeBubbleStartTime().Add(48 * time.Hour)),
		}
	}
	ProjectMember1 = func() *apiv1.ProjectMember {
		return &apiv1.ProjectMember{
			Id:                  "16d6e8ba-f574-494f-8d5e-74f6cb2d8db0",
			Role:                apiv1.ProjectRole_PROJECT_ROLE_OWNER,
			InheritedMembership: false,
			CreatedAt:           timestamppb.New(e2e.TimeBubbleStartTime()),
		}
	}
	ProjectMember2 = func() *apiv1.ProjectMember {
		return &apiv1.ProjectMember{
			Id:                  "40c0da4b-9eb9-4371-91aa-1ae62193fa54",
			Role:                apiv1.ProjectRole_PROJECT_ROLE_EDITOR,
			InheritedMembership: true,
			CreatedAt:           timestamppb.New(e2e.TimeBubbleStartTime()),
		}
	}
)
