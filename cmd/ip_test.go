package cmd

import (
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/google/go-cmp/cmp/cmpopts"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func Test_IPCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.IP]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.IP) []string {
				return []string{"ip", "list", "--project", "a"}
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("List", mock.Anything, connect.NewRequest(&apiv1.IPServiceListRequest{
						Project: "a",
					})).Return(&connect.Response[apiv1.IPServiceListResponse]{
						Msg: &apiv1.IPServiceListResponse{
							Ips: []*apiv1.IP{
								{
									Uuid:        "uuid",
									Ip:          "4.3.2.1",
									Name:        "b",
									Description: "b description",
									Network:     "b-network",
									Project:     "b",
									Type:        "ephemeral",
									Tags:        []string{"a=b"},
								},
								{
									Uuid:        "uuid",
									Ip:          "1.1.1.1",
									Name:        "a",
									Description: "a description",
									Network:     "a-network",
									Project:     "a",
									Type:        "ephemeral",
									Tags:        []string{"a=b"},
								},
							},
						},
					}, nil)
				},
			},
			Want: []*apiv1.IP{
				{
					Uuid:        "uuid",
					Ip:          "1.1.1.1",
					Name:        "a",
					Description: "a description",
					Network:     "a-network",
					Project:     "a",
					Type:        "ephemeral",
					Tags:        []string{"a=b"},
				},
				{
					Uuid:        "uuid",
					Ip:          "4.3.2.1",
					Name:        "b",
					Description: "b description",
					Network:     "b-network",
					Project:     "b",
					Type:        "ephemeral",
					Tags:        []string{"a=b"},
				},
			},
			WantTable: pointer.Pointer(`
IP        PROJECT
1.1.1.1   a
4.3.2.1   b
`),
			WantWideTable: pointer.Pointer(`
IP        PROJECT
1.1.1.1   a
4.3.2.1   b
`),
			Template: pointer.Pointer("{{ .ip }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
1.1.1.1 a
4.3.2.1 b
			`),
			WantMarkdown: pointer.Pointer(`
|   IP    | PROJECT |
|---------|---------|
| 1.1.1.1 | a       |
| 4.3.2.1 | b       |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_SingleResult(t *testing.T) {
	tests := []*Test[*apiv1.IP]{
		{
			Name: "describe",
			Cmd: func(want *apiv1.IP) []string {
				return []string{"ip", "describe", "--project", "a", "uuid"}
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceGetRequest{
						Project: "a",
						Uuid:    "uuid",
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceGetResponse]{
						Msg: &apiv1.IPServiceGetResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Project:     "a",
								Type:        "ephemeral",
								Tags:        []string{"a=b"},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        "ephemeral",
				Tags:        []string{"a=b"},
			},
			WantTable: pointer.Pointer(`
IP        PROJECT
1.1.1.1   a
`),
			WantWideTable: pointer.Pointer(`
IP        PROJECT
1.1.1.1   a
`),
			Template: pointer.Pointer("{{ .ip }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
1.1.1.1 a
			`),
			WantMarkdown: pointer.Pointer(`
|   IP    | PROJECT |
|---------|---------|
| 1.1.1.1 | a       |
`),
		},
		{
			Name: "delete",
			Cmd: func(want *apiv1.IP) []string {
				return []string{"ip", "rm", "--project", "a", "uuid"}
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
						Project: "a",
						Uuid:    "uuid",
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceDeleteResponse]{
						Msg: &apiv1.IPServiceDeleteResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Project:     "a",
								Type:        "ephemeral",
								Tags:        []string{"a=b"},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        "ephemeral",
				Tags:        []string{"a=b"},
			},
		},
		{
			Name: "create",
			Cmd: func(want *apiv1.IP) []string {
				args := []string{"ip", "create", "--project", "a", "--description", "a description", "--name", "a", "--network", "a-network", "--static", "--tags", "a=b"}
				AssertExhaustiveArgs(t, args, "file")
				return args
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
						Project:     "a",
						Name:        "a",
						Description: "a description",
						Network:     "a-network",
						Static:      true,
						Tags:        []string{"a=b"},
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceAllocateResponse]{
						Msg: &apiv1.IPServiceAllocateResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Project:     "a",
								Type:        "ephemeral",
								Tags:        []string{"a=b"},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        "ephemeral",
				Tags:        []string{"a=b"},
			},
		},
		{
			Name: "update",
			Cmd: func(want *apiv1.IP) []string {
				args := []string{"ip", "update", "--project", "a", "--uuid", "uuid", "--description", "b description", "--name", "b", "--static", "--tags", "c=d"}
				AssertExhaustiveArgs(t, args, "file")
				return args
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceGetRequest{
						Uuid:    "uuid",
						Project: "a",
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceGetResponse]{
						Msg: &apiv1.IPServiceGetResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Project:     "a",
								Type:        "ephemeral",
								Tags:        []string{"a=b"},
							},
						},
					}, nil)
					m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
						Ip: &apiv1.IP{
							Uuid:        "uuid",
							Ip:          "1.1.1.1",
							Name:        "b",
							Description: "b description",
							Network:     "a-network",
							Project:     "a",
							Type:        "static",
							Tags:        []string{"c=d"},
						},
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceUpdateResponse]{
						Msg: &apiv1.IPServiceUpdateResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "b",
								Description: "b description",
								Network:     "a-network",
								Project:     "a",
								Type:        "static",
								Tags:        []string{"c=d"},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "b",
				Description: "b description",
				Network:     "a-network",
				Project:     "a",
				Type:        "static",
				Tags:        []string{"c=d"},
			},
		},
		{
			Name: "update from file",
			Cmd: func(want *apiv1.IP) []string {
				return []string{"ip", "update", "-f", "/file.yaml"}
			},
			FsMocks: func(fs afero.Fs, want *apiv1.IP) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshal(t, want), 0755))
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
						Ip: &apiv1.IP{
							Uuid:        "uuid",
							Ip:          "1.1.1.1",
							Project:     "a",
							Name:        "a",
							Description: "a description",
							Network:     "a-network",
							Type:        "static",
							Tags:        []string{"a=b"},
						},
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceUpdateResponse]{
						Msg: &apiv1.IPServiceUpdateResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Project:     "a",
								Type:        "static",
								Tags:        []string{"a=b"},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        "static",
				Tags:        []string{"a=b"},
			},
		},
		{
			Name: "create from file",
			Cmd: func(want *apiv1.IP) []string {
				return []string{"ip", "create", "-f", "/file.yaml"}
			},
			FsMocks: func(fs afero.Fs, want *apiv1.IP) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshal(t, want), 0755))
			},
			APIMocks: &apitests.Apiv1MockFns{
				IP: func(m *mock.Mock) {
					m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
						Project:     "a",
						Name:        "a",
						Description: "a description",
						Network:     "a-network",
						Static:      true,
						Tags:        []string{"a=b"},
					}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(&connect.Response[apiv1.IPServiceAllocateResponse]{
						Msg: &apiv1.IPServiceAllocateResponse{
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Project:     "a",
								Type:        "static",
								Tags:        []string{"a=b"},
							},
						},
					}, nil)
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        "static",
				Tags:        []string{"a=b"},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
