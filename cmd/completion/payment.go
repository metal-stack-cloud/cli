package completion

import (
	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	"github.com/spf13/cobra"
)

func (c *Completion) AdminPaymentCouponListCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	req := &adminv1.PaymentServiceListCouponsRequest{}
	resp, err := c.Client.Adminv1().Payment().ListCoupons(c.Ctx, connect.NewRequest(req))
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, s := range resp.Msg.Coupons {
		names = append(names, s.Id+"\t"+s.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
