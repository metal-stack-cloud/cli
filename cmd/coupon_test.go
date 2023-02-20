package cmd

import (
	"testing"

	"github.com/bufbuild/connect-go"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/stretchr/testify/mock"
)

func Test_CouponCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.Coupon] {
		{
			Name: "list coupons",
			Cmd: func(want []*apiv1.Coupon) []string {
				return []string{"admin", "coupon", "list"}
			},
			AdminMocks: &apitests.Adminv1MockFns{
				Payment: func(m *mock.Mock) {
					m.On("ListCoupons", mock.Anything, connect.NewRequest(&adminv1.PaymentServiceListCouponsRequest{})).Return(&connect.Response[adminv1.PaymentServiceListCouponsResponse]{
						Msg: &adminv1.PaymentServiceListCouponsResponse{
							Coupons: []*apiv1.Coupon{
								{
									Id: "someidone",
									Name: "somenameone",
									AmountOff: 150,
									DurationInMonth: 4,
									TimesRedeemed: 1,
									MaxRedemptions: 30,
									Currency: "eur",
									// CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
								},
								{
									Id: "someidtwo",
									Name: "somenametwo",
									AmountOff: 50,
									DurationInMonth: 3,
									TimesRedeemed: 2,
									MaxRedemptions: 150,
									Currency: "eur",
									// CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
								},
							},
						},
					}, nil)
				},
			},
			Want: []*apiv1.Coupon{
				{
					Id: "someidone",
					Name: "somenameone",
					AmountOff: 150,
					DurationInMonth: 4,
					TimesRedeemed: 1,
					MaxRedemptions: 30,
					Currency: "eur",
					// CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
				},
				{
					Id: "someidtwo",
					Name: "somenametwo",
					AmountOff: 50,
					DurationInMonth: 3,
					TimesRedeemed: 2,
					MaxRedemptions: 150,
					Currency: "eur",
					// CreatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
				},
			},
// 			WantTable: pointer.Pointer(`
// ID        NAME        AMOUNTOFF DURATION REDEEMED CREATED
// someidone somenameone 1.00      4 month  0/30     1970-01-01 00:00:00 +0000 UTC
// someidtwo somenametwo 0.00      3 month  0/150    1970-01-01 00:00:00 +0000 UTC
// `),
// 			WantWideTable: pointer.Pointer(`
// ID         NAME            AMOUNTOFF    DURATION   REDEEMED   CREATED
// someidone  somenameone     150          4          0/30       1970-01-01 00:00:00 +0000 UTC
// someidtwo  somenametwo     50           3          0/150      1970-01-01 00:00:00 +0000 UTC
// `),
// 			Template: pointer.Pointer("{{ .id }} {{ .name }} {{ .amount_off }} {{ .duration_in_month }} {{ .times_redeemed }}/{{ .max_redemptions }} {{ .created_at }}"),
// 			WantTemplate: pointer.Pointer(`
// someidone somenameone 150 4 1/30 1970-01-01 00:00:00 +0000 UTC
// someidtwo somenametwo 50 3 2/150 1970-01-01 00:00:00 +0000 UTC
// `),
// 			WantMarkdown: pointer.Pointer(`
// | ID        | NAME        | AMOUNTOFF  | DURATION | REDEEMED | CREATED                       |
// |-----------|-------------|------------|----------|----------|-------------------------------|
// | someidone | somenameone | 150.00 eur | 4 month  | 0/30     | 1970-01-01 00:00:00 +0000 UTC |
// | someidtwo | somenametwo | 50.00 eur  | 3 month  | 0/150    | 1970-01-01 00:00:00 +0000 UTC |
// `),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}