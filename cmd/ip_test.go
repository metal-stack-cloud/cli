package cmd

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp/cmpopts"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	v1 "github.com/metal-stack-cloud/cli/cmd/api/v1"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/metal-stack/metal-lib/pkg/tag"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
)

var (
	ip1 = func() *apiv1.IP {
		return &apiv1.IP{
			Uuid:        "2e0144a2-09ef-42b7-b629-4263295db6e8",
			Ip:          "1.1.1.1",
			Name:        "a",
			Description: "a description",
			Project:     "a",
			Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
			Tags:        []string{tag.New(tag.ClusterServiceFQN, "<cluster>/default/ingress-nginx")},
		}
	}
	ip2 = func() *apiv1.IP {
		return &apiv1.IP{
			Uuid:        "9cef40ec-29c6-4dfa-aee8-47ee1f49223d",
			Ip:          "4.3.2.1",
			Name:        "b",
			Description: "b description",
			Project:     "b",
			Type:        apiv1.IPType_IP_TYPE_EPHEMERAL,
			Tags:        []string{"a=b"},
		}
	}
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
									ip2(),
									ip1(),
								},
							},
						}, nil)
					},
				},
			},
			Want: []*apiv1.IP{
				ip1(),
				ip2(),
			},
			WantTable: pointer.Pointer(`
IP        PROJECT   ID                                     TYPE        NAME   ATTACHED SERVICE
1.1.1.1   a         2e0144a2-09ef-42b7-b629-4263295db6e8   ephemeral   a      ingress-nginx
4.3.2.1   b         9cef40ec-29c6-4dfa-aee8-47ee1f49223d   ephemeral   b
`),
			WantWideTable: pointer.Pointer(`
IP        PROJECT   ID                                     TYPE        NAME   DESCRIPTION     LABELS
1.1.1.1   a         2e0144a2-09ef-42b7-b629-4263295db6e8   ephemeral   a      a description   cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
4.3.2.1   b         9cef40ec-29c6-4dfa-aee8-47ee1f49223d   ephemeral   b      b description   a=b
`),
			Template: pointer.Pointer("{{ .ip }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
1.1.1.1 a
4.3.2.1 b
			`),
			WantMarkdown: pointer.Pointer(`
|   IP    | PROJECT |                  ID                  |   TYPE    | NAME | ATTACHED SERVICE |
|---------|---------|--------------------------------------|-----------|------|------------------|
| 1.1.1.1 | a       | 2e0144a2-09ef-42b7-b629-4263295db6e8 | ephemeral | a    | ingress-nginx    |
| 4.3.2.1 | b       | 9cef40ec-29c6-4dfa-aee8-47ee1f49223d | ephemeral | b    |                  |
`),
		},
		{
			Name: "apply",
			Cmd: func(want []*apiv1.IP) []string {
				return appendFromFileCommonArgs("ip", "apply")
			},
			FsMocks: func(fs afero.Fs, want []*apiv1.IP) {
				require.NoError(t, afero.WriteFile(fs, "/file.yaml", MustMarshalToMultiYAML(t, want), 0755))
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToCreate(ip1())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
							Ip: ip1(),
						}), nil)
						m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(v1.IpResponseToCreate(ip2())), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
							Ip: ip2(),
						}), nil)
						// FIXME: API does not return a conflict when already exists, so this functionality does not work!
					},
				},
			},
			Want: []*apiv1.IP{
				ip1(),
				ip2(),
			},
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
				return []string{"ip", "describe", "--project", want.Project, want.Uuid}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceGetRequest{
							Project: ip1().Project,
							Uuid:    ip1().Uuid,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceGetResponse{
							Ip: ip1(),
						}), nil)
					},
				},
			},
			Want: ip1(),
			WantTable: pointer.Pointer(`
IP        PROJECT   ID                                     TYPE        NAME   ATTACHED SERVICE
1.1.1.1   a         2e0144a2-09ef-42b7-b629-4263295db6e8   ephemeral   a      ingress-nginx
`),
			WantWideTable: pointer.Pointer(`
IP        PROJECT   ID                                     TYPE        NAME   DESCRIPTION     LABELS
1.1.1.1   a         2e0144a2-09ef-42b7-b629-4263295db6e8   ephemeral   a      a description   cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
`),
			Template: pointer.Pointer("{{ .ip }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
1.1.1.1 a
			`),
			WantMarkdown: pointer.Pointer(`
|   IP    | PROJECT |                  ID                  |   TYPE    | NAME | ATTACHED SERVICE |
|---------|---------|--------------------------------------|-----------|------|------------------|
| 1.1.1.1 | a       | 2e0144a2-09ef-42b7-b629-4263295db6e8 | ephemeral | a    | ingress-nginx    |
`),
		},
		{
			Name: "delete",
			Cmd: func(want *apiv1.IP) []string {
				return []string{"ip", "rm", "--project", want.Project, want.Uuid}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Delete", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
							Project: ip1().Project,
							Uuid:    ip1().Uuid,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceDeleteResponse{
							Ip: ip1(),
						}), nil)
					},
				},
			},
			Want: ip1(),
		},
		{
			Name: "create",
			Cmd: func(want *apiv1.IP) []string {
				args := []string{"ip", "create", "--project", want.Project, "--description", want.Description, "--name", want.Name, "--static", "--tags", "a=b"}
				AssertExhaustiveArgs(t, args, commonExcludedFileArgs()...)
				return args
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Allocate", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
							Project:     ip1().Project,
							Name:        ip1().Name,
							Description: ip1().Description,
							Static:      true,
							Tags:        []string{"a=b"},
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
							Ip: ip1(),
						}), nil)
					},
				},
			},
			Want: ip1(),
		},
		{
			Name: "update",
			Cmd: func(want *apiv1.IP) []string {
				args := []string{"ip", "update", "--project", want.Project, want.Uuid, "--description", "b description", "--name", "b", "--static", "--tags", "c=d"}
				AssertExhaustiveArgs(t, args, commonExcludedFileArgs()...)
				return args
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					IP: func(m *mock.Mock) {
						m.On("Get", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceGetRequest{
							Uuid:    ip1().Uuid,
							Project: ip1().Project,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceGetResponse{
							Ip: ip1(),
						}), nil)

						m.On("Update", mock.Anything, testcommon.MatchByCmpDiff(t, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
							Project: ip1().Project,
							Ip: &apiv1.IP{
								Uuid:        ip1().Uuid,
								Ip:          "1.1.1.1",
								Name:        "b",
								Description: "b description",
								Project:     ip1().Project,
								Type:        apiv1.IPType_IP_TYPE_STATIC,
								Tags:        []string{"c=d"},
							},
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceUpdateResponse{
							Ip: ip1(),
						}), nil)
					},
				},
			},
			Want: ip1(),
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
							Project: ip1().Project,
							Ip:      ip1(),
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceUpdateResponse{
							Ip: ip1(),
						}), nil)
					},
				},
			},
			Want: ip1(),
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
							Project:     ip1().Project,
							Name:        ip1().Name,
							Description: ip1().Description,
							Static:      false,
							Tags:        ip1().Tags,
						}), cmpopts.IgnoreTypes(protoimpl.MessageState{}))).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
							Ip: ip1(),
						}), nil)
					},
				},
			},
			Want: ip1(),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
