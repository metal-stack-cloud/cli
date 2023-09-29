package cmd

import (
	"testing"

	"connectrpc.com/connect"
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
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
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
										Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
										Tags:        []string{"a=b"},
									},
									{
										Uuid:        "uuid",
										Ip:          "1.1.1.1",
										Name:        "a",
										Description: "a description",
										Network:     "a-network",
										Project:     "a",
										Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
										Tags:        []string{"a=b"},
									},
								},
							},
						}, nil)
					},
				}},

			Want: []*apiv1.IP{
				{
					Uuid:        "uuid",
					Ip:          "1.1.1.1",
					Name:        "a",
					Description: "a description",
					Network:     "a-network",
					Project:     "a",
					Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
					Tags:        []string{"a=b"},
				},
				{
					Uuid:        "uuid",
					Ip:          "4.3.2.1",
					Name:        "b",
					Description: "b description",
					Network:     "b-network",
					Project:     "b",
					Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
					Tags:        []string{"a=b"},
				},
			},
			WantTable: pointer.Pointer(`
IP        ID     PROJECT   NAME   DESCRIPTION     TYPE
1.1.1.1   uuid   a         a      a description   ephemeral
4.3.2.1   uuid   b         b      b description   ephemeral
`),
			WantWideTable: pointer.Pointer(`
IP        ID     PROJECT   NAME   DESCRIPTION     TYPE
1.1.1.1   uuid   a         a      a description   ephemeral
4.3.2.1   uuid   b         b      b description   ephemeral
`),
			Template: pointer.Pointer("{{ .ip }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
1.1.1.1 a
4.3.2.1 b
			`),
			WantMarkdown: pointer.Pointer(`
|   IP    |  ID  | PROJECT | NAME |  DESCRIPTION  |   TYPE    |
|---------|------|---------|------|---------------|-----------|
| 1.1.1.1 | uuid | a       | a    | a description | ephemeral |
| 4.3.2.1 | uuid | b       | b    | b description | ephemeral |
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
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
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
									Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
									Tags:        []string{"a=b"},
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
				Tags:        []string{"a=b"},
			},
			WantTable: pointer.Pointer(`
IP        ID     PROJECT   NAME   DESCRIPTION     TYPE
1.1.1.1   uuid   a         a      a description   ephemeral
`),
			WantWideTable: pointer.Pointer(`
IP        ID     PROJECT   NAME   DESCRIPTION     TYPE
1.1.1.1   uuid   a         a      a description   ephemeral
`),
			Template: pointer.Pointer("{{ .ip }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
1.1.1.1 a
			`),
			WantMarkdown: pointer.Pointer(`
|   IP    |  ID  | PROJECT | NAME |  DESCRIPTION  |   TYPE    |
|---------|------|---------|------|---------------|-----------|
| 1.1.1.1 | uuid | a       | a    | a description | ephemeral |
`),
		},
		{
			Name: "delete",
			Cmd: func(want *apiv1.IP) []string {
				return []string{"ip", "rm", "--project", "a", "uuid"}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
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
									Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
									Tags:        []string{"a=b"},
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
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
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
							Project:     "a",
							Name:        "a",
							Description: "a description",
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
									Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
									Tags:        []string{"a=b"},
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
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
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
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
									Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
									Tags:        []string{"a=b"},
								},
							},
						}, nil)
						m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
							Project: "a",
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Name:        "b",
								Description: "b description",
								Network:     "a-network",
								Project:     "a",
								Type:        apiv1.IPType_IP_TYPE_STATIC,
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
									Type:        apiv1.IPType_IP_TYPE_STATIC,
									Tags:        []string{"c=d"},
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "b",
				Description: "b description",
				Network:     "a-network",
				Project:     "a",
				Type:        apiv1.IPType_IP_TYPE_STATIC,
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
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
							Project: "a",
							Ip: &apiv1.IP{
								Uuid:        "uuid",
								Ip:          "1.1.1.1",
								Project:     "a",
								Name:        "a",
								Description: "a description",
								Network:     "a-network",
								Type:        apiv1.IPType_IP_TYPE_STATIC,
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
									Type:        apiv1.IPType_IP_TYPE_STATIC,
									Tags:        []string{"a=b"},
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        apiv1.IPType_IP_TYPE_STATIC,
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
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
							Project:     "a",
							Name:        "a",
							Description: "a description",
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
									Type:        apiv1.IPType_IP_TYPE_STATIC,
									Tags:        []string{"a=b"},
								},
							},
						}, nil)
					},
				},
			},
			Want: &apiv1.IP{
				Uuid:        "uuid",
				Ip:          "1.1.1.1",
				Name:        "a",
				Description: "a description",
				Network:     "a-network",
				Project:     "a",
				Type:        apiv1.IPType_IP_TYPE_STATIC,
				Tags:        []string{"a=b"},
			},
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
