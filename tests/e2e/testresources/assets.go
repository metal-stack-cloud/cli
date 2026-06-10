package testresources

import apiv1 "github.com/metal-stack-cloud/api/go/api/v1"

var (
	Asset1 = func() *apiv1.Asset {
		return &apiv1.Asset{
			Region: &apiv1.Region{
				Id:     "fra",
				Name:   "Frankfurt",
				Active: true,
				Partitions: map[string]*apiv1.Partition{
					"fra-1": &apiv1.Partition{
						Id: "fra-1",
					},
				},
			},
			MachineTypes: map[string]*apiv1.MachineType{
				"c1-medium-x86": &apiv1.MachineType{
					Id: "c1-medium-x86",
				},
			},
		}
	}
)
