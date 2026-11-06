package testresources

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/tag"
)

var (
	Ip1 = func() *apiv1.IP {
		return &apiv1.IP{
			Uuid:        "2e0144a2-09ef-42b7-b629-4263295db6e8",
			Ip:          "1.1.1.1",
			Name:        "ip1",
			Description: "ip1 description",
			Project:     Project1().Uuid,
			Type:        apiv1.IPType_IP_TYPE_STATIC,
			Tags:        []string{tag.New(tag.ClusterServiceFQN, "<cluster>/default/ingress-nginx")},
		}
	}
	Ip2 = func() *apiv1.IP {
		return &apiv1.IP{
			Uuid:        "9cef40ec-29c6-4dfa-aee8-47ee1f49223d",
			Ip:          "4.3.2.1",
			Name:        "ip2",
			Description: "ip2 description",
			Project:     Project2().Uuid,
			Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
			Tags:        []string{"a=b"},
		}
	}
)
