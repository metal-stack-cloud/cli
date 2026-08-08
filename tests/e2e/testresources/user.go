package testresources

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

var (
	User1 = func() *apiv1.User {
		return &apiv1.User{
			Login:          Tenant1().Login,
			Name:           "Larry",
			Email:          "evilLarry@mail.com",
			OauthProvider:  apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
			DefaultTenant:  Tenant1(),
			DefaultProject: Project1(),
			Tenants:        []*apiv1.Tenant{Tenant1(), Tenant2()},
			Projects:       []*apiv1.Project{Project1(), Project2()},
			AvatarUrl:      "https://tenor.com/de/view/evil-larry-larry-gif-431258253732208458",
		}
	}
	User2 = func() *apiv1.User {
		return &apiv1.User{
			Login:          Tenant2().Login,
			Name:           "Timmy",
			Email:          "timmy@mail.com",
			OauthProvider:  apiv1.OAuthProvider_O_AUTH_PROVIDER_GITHUB,
			DefaultTenant:  Tenant2(),
			DefaultProject: Project2(),
			Tenants:        []*apiv1.Tenant{Tenant2()},
			Projects:       []*apiv1.Project{Project2()},
		}
	}
)
