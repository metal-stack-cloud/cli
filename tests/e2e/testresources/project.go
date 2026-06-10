package testresources

import (
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
)
