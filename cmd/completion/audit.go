package completion

import (
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) AuditPhaseListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{apiv1.AuditPhase_AUDIT_PHASE_REQUEST.String(), apiv1.AuditPhase_AUDIT_PHASE_RESPONSE.String()}, cobra.ShellCompDirectiveNoFileComp
}
