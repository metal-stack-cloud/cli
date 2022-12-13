package cmd

import (
	"testing"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func Test_FirewallCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.IP]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.IP) []string {
				return []string{"ip", "list"}
			},
			Want: []*apiv1.IP{
				{
					Ip: "1.1.1.1",
				},
			},
			WantTable: pointer.Pointer(`

`),
			WantWideTable: pointer.Pointer(`
ID   AGE   HOSTNAME              PROJECT     NETWORKS   IPS       PARTITION
1    14d   firewall-hostname-1   project-1   private    1.1.1.1   1
2    14d   firewall-hostname-2   project-1   private    1.1.1.1   1
`),
			// Template: pointer.Pointer("{{ .id }} {{ .name }}"),
			// 			WantTemplate: pointer.Pointer(`
			// 1 firewall-1
			// 2 firewall-2
			// `),
			WantMarkdown: pointer.Pointer(`
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
