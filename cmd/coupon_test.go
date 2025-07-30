package cmd

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_CouponCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.Coupon]{
		{
			Name: "list coupons",
			Cmd: func(want []*apiv1.Coupon) []string {
				return []string{"admin", "coupon", "list"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Adminv1Mocks: &apitests.Adminv1MockFns{
					Payment: func(m *mock.Mock) {
						m.On("ListCoupons", mock.Anything, connect.NewRequest(&adminv1.PaymentServiceListCouponsRequest{})).Return(&connect.Response[adminv1.PaymentServiceListCouponsResponse]{
							Msg: &adminv1.PaymentServiceListCouponsResponse{
								Coupons: []*apiv1.Coupon{
									{
										Id:              "someidone",
										Name:            "somenameone",
										AmountOff:       150,
										DurationInMonth: 4,
										TimesRedeemed:   1,
										MaxRedemptions:  30,
										Currency:        "eur",
										CreatedAt:       timestamppb.New(time.Now().Add(-1 * time.Hour)),
									},
									{
										Id:              "someidtwo",
										Name:            "somenametwo",
										AmountOff:       50,
										DurationInMonth: 3,
										TimesRedeemed:   2,
										MaxRedemptions:  150,
										Currency:        "eur",
										CreatedAt:       timestamppb.New(time.Now().Add(-1 * time.Hour)),
									},
								},
							},
						}, nil)
					},
				},
			},
			Want: []*apiv1.Coupon{
				{
					Id:              "someidone",
					Name:            "somenameone",
					AmountOff:       150,
					DurationInMonth: 4,
					TimesRedeemed:   1,
					MaxRedemptions:  30,
					Currency:        "eur",
					CreatedAt:       timestamppb.New(time.Now().Add(-1 * time.Hour)),
				},
				{
					Id:              "someidtwo",
					Name:            "somenametwo",
					AmountOff:       50,
					DurationInMonth: 3,
					TimesRedeemed:   2,
					MaxRedemptions:  150,
					Currency:        "eur",
					CreatedAt:       timestamppb.New(time.Now().Add(-1 * time.Hour)),
				},
			},
			WantTable: pointer.Pointer(`
ID         NAME         AMOUNT OFF  DURATION  REDEEMED  CREATED     
someidone  somenameone  1.00 eur    4 month   1/30      1 hour ago  
someidtwo  somenametwo  0.00 eur    3 month   2/150     1 hour ago
`),
			WantWideTable: pointer.Pointer(`
ID         NAME         AMOUNT OFF  DURATION  REDEEMED  CREATED     
someidone  somenameone  1.00 eur    4 month   1/30      1 hour ago  
someidtwo  somenametwo  0.00 eur    3 month   2/150     1 hour ago
`),
			Template: pointer.Pointer(`{{ .id }} {{ .name }} {{ .amount_off }} {{ .duration_in_month }} {{ .times_redeemed }}/{{ .max_redemptions }} {{ date "02/01/2006" .created_at }}`),
			WantTemplate: pointer.Pointer(`
someidone somenameone 150 4 1/30 19/05/2022
someidtwo somenametwo 50 3 2/150 19/05/2022
			`),
			WantMarkdown: pointer.Pointer(`
| ID        | NAME        | AMOUNT OFF | DURATION | REDEEMED | CREATED    |
|-----------|-------------|------------|----------|----------|------------|
| someidone | somenameone | 1.00 eur   | 4 month  | 1/30     | 1 hour ago |
| someidtwo | somenametwo | 0.00 eur   | 3 month  | 2/150    | 1 hour ago |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
