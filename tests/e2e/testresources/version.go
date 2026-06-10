package testresources

import apiv1 "github.com/metal-stack-cloud/api/go/api/v1"

var (
	Version = func() *apiv1.Version {
		return &apiv1.Version{
			Version:   "v0.1.8",
			Revision:  "tags/v0.1.8-0-g476edc0",
			GitSha1:   "477edc0b",
			BuildDate: "2026-03-21T15:35:07+00:00",
		}
	}
)
