package api_e2e

import (
	"strings"
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	e2erootcmd "github.com/metal-stack-cloud/cli/testing/e2e"
	"github.com/metal-stack-cloud/cli/tests/e2e/testresources"
	e2e "github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_IPCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.IPServiceListResponse, []*apiv1.IP]{
		{
			Name:    "list",
			CmdArgs: []string{"ip", "list", "--project", testresources.Project1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						IP: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.IPServiceListRequest{
								Project: testresources.Project1().Uuid,
							})).Return(connect.NewResponse(&apiv1.IPServiceListResponse{
								Ips: []*apiv1.IP{
									testresources.Ip1(),
								},
							}), nil)
						},
					},
				},
			}),

			WantTable: new(`
            IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE  
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ingress-nginx
`),
			WantWideTable: new(`
            IP       PROJECT                               ID                                    TYPE    NAME  DESCRIPTION      LABELS                                                                       
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ip1 description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
`),
			Template:     new("{{ .ip }} {{ .project }}"),
			WantTemplate: new(`1.1.1.1 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | IP      | PROJECT                              | ID                                   | TYPE   | NAME | ATTACHED SERVICE |
            |---------|--------------------------------------|--------------------------------------|--------|------|------------------|
            | 1.1.1.1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static | ip1  | ingress-nginx    |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Apply(t *testing.T) {
	tests := []*e2e.Test[apiv1.IPServiceListResponse, []*apiv1.IP]{
		{
			Name:    "apply ip",
			CmdArgs: append([]string{"ip", "apply"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Ip1(), testresources.Ip2()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							IP: func(m *mock.Mock) {
								m.On("Allocate", mock.Anything, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
									Project:     testresources.Ip2().Project,
									Name:        testresources.Ip2().Name,
									Description: testresources.Ip2().Description,
									Tags:        testresources.Ip2().Tags,
								})).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
									Ip: testresources.Ip2(),
								}), nil)
								m.On("Allocate", mock.Anything, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
									Project:     testresources.Ip1().Project,
									Name:        testresources.Ip1().Name,
									Description: testresources.Ip1().Description,
									Tags:        testresources.Ip1().Tags,
									Static:      true,
								})).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
									Ip: testresources.Ip1(),
								}), nil)
								// FIXME: API does not return a conflict when already exists, so the update functionality does not work!
							},
						},
					},
				}),

			WantTable: new(`
			IP       PROJECT                               ID                                    TYPE       NAME  ATTACHED SERVICE  
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static     ip1   ingress-nginx     
            4.3.2.1  c40ad996-e1fd-4511-a7bf-418219cb8d67  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  ip2
			`),
			WantWideTable: new(`
            IP       PROJECT                               ID                                    TYPE       NAME  DESCRIPTION      LABELS                                                                       
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static     ip1   ip1 description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx  
            4.3.2.1  c40ad996-e1fd-4511-a7bf-418219cb8d67  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  ip2   ip2 description  a=b
			`),
			WantMarkdown: new(`
            | IP      | PROJECT                              | ID                                   | TYPE      | NAME | ATTACHED SERVICE |
            |---------|--------------------------------------|--------------------------------------|-----------|------|------------------|
            | 1.1.1.1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static    | ip1  | ingress-nginx    |
            | 4.3.2.1 | c40ad996-e1fd-4511-a7bf-418219cb8d67 | 9cef40ec-29c6-4dfa-aee8-47ee1f49223d | ephemeral | ip2  |                  |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.IPServiceUpdateResponse, *apiv1.IP]{
		{
			Name: "update ip",
			CmdArgs: []string{"ip", "update", testresources.Ip1().Uuid,
				"--project", testresources.Project1().Uuid,
				"--description", testresources.Ip1().Description,
				"--name", testresources.Ip1().Name,
				"--tags", strings.Join(testresources.Ip1().Tags, ","),
				"--static=true"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						IP: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.IPServiceGetRequest{
								Uuid:    testresources.Ip1().Uuid,
								Project: testresources.Project1().Uuid,
							})).Return(connect.NewResponse(&apiv1.IPServiceGetResponse{
								Ip: testresources.Ip1(),
							}), nil)
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
								Project: testresources.Project1().Uuid,
								Ip:      testresources.Ip1(),
							})).Return(connect.NewResponse(&apiv1.IPServiceUpdateResponse{
								Ip: testresources.Ip1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Ip1(),
			WantDefault: new(`
description: ip1 description
ip: 1.1.1.1
name: ip1
project: c40ad996-e1fd-4511-a7bf-418219cb8d95
tags:
- cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
type: IP_TYPE_STATIC
uuid: 2e0144a2-09ef-42b7-b629-4263295db6e8
			`),
		},
		{
			Name:    "update ip from file",
			CmdArgs: append([]string{"ip", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Ip1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							IP: func(m *mock.Mock) {
								m.On("Update", mock.Anything, connect.NewRequest(&apiv1.IPServiceUpdateRequest{
									Project: testresources.Project1().Uuid,
									Ip:      testresources.Ip1(),
								})).Return(connect.NewResponse(&apiv1.IPServiceUpdateResponse{
									Ip: testresources.Ip1(),
								}), nil)
							},
						},
					},
				}),
			WantTable: new(`
            IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE  
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ingress-nginx
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Create(t *testing.T) {
	tests := []*e2e.Test[apiv1.IPServiceAllocateResponse, *apiv1.IP]{
		{
			Name: "create",
			CmdArgs: []string{"ip", "create",
				"--project", testresources.Project1().Uuid,
				"--description", testresources.Ip1().Description,
				"--name", testresources.Ip1().Name,
				"--tags", strings.Join(testresources.Ip1().Tags, ","),
				"--static=true"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						IP: func(m *mock.Mock) {
							m.On("Allocate", mock.Anything, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
								Project:     testresources.Project1().Uuid,
								Name:        testresources.Ip1().Name,
								Description: testresources.Ip1().Description,
								Tags:        testresources.Ip1().Tags,
								Static:      true,
							})).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
								Ip: testresources.Ip1(),
							}), nil)
						},
					},
				},
			}),

			WantObject: testresources.Ip1(),
		},
		{
			Name:    "create from file",
			CmdArgs: append([]string{"ip", "create"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Ip1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							IP: func(m *mock.Mock) {
								m.On("Allocate", mock.Anything, connect.NewRequest(&apiv1.IPServiceAllocateRequest{
									Project:     testresources.Ip1().Project,
									Name:        testresources.Ip1().Name,
									Description: testresources.Ip1().Description,
									Tags:        testresources.Ip1().Tags,
									Static:      true,
								})).Return(connect.NewResponse(&apiv1.IPServiceAllocateResponse{
									Ip: testresources.Ip1(),
								}), nil)
							},
						},
					},
				}),

			WantTable: new(`
			IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE  
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ingress-nginx
			`),
			WantWideTable: new(`
			IP       PROJECT                               ID                                    TYPE    NAME  DESCRIPTION      LABELS                                                                       
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ip1 description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.IPServiceDeleteResponse, *apiv1.IP]{
		{
			Name:    "delete ip",
			CmdArgs: []string{"ip", "rm", "--project", testresources.Project1().Uuid, testresources.Ip1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						IP: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
								Project: testresources.Project1().Uuid,
								Uuid:    testresources.Ip1().Uuid,
							})).Return(connect.NewResponse(&apiv1.IPServiceDeleteResponse{
								Ip: testresources.Ip1(),
							}), nil)
						},
					},
				}}),

			WantObject: testresources.Ip1(),
		},
		{
			Name:    "delete ip from file",
			CmdArgs: append([]string{"ip", "delete"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t,
				&e2erootcmd.TestConfig{
					FsMocks: func(fs *afero.Afero) {
						require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Ip1()), 0755))
					},
					ClientMocks: &apitests.ClientMockFns{
						Apiv1Mocks: &apitests.Apiv1MockFns{
							IP: func(m *mock.Mock) {
								m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.IPServiceDeleteRequest{
									Project: testresources.Project1().Uuid,
									Uuid:    testresources.Ip1().Uuid,
								})).Return(connect.NewResponse(&apiv1.IPServiceDeleteResponse{
									Ip: testresources.Ip1(),
								}), nil)
							},
						},
					},
				}),

			WantDefault: new(`
IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE  
1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ingress-nginx     
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.IPServiceGetResponse, *apiv1.IP]{
		{
			Name:    "describe ip",
			CmdArgs: []string{"ip", "describe", "--project", testresources.Project1().Uuid, testresources.Ip1().Uuid},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						IP: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.IPServiceGetRequest{
								Project: testresources.Project1().Uuid,
								Uuid:    testresources.Ip1().Uuid,
							})).Return(connect.NewResponse(&apiv1.IPServiceGetResponse{
								Ip: testresources.Ip1(),
							}), nil)
						},
					},
				}}),

			WantObject: testresources.Ip1(),
			WantTable: new(`
            IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE  
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ingress-nginx
`),
			WantWideTable: new(`
            IP       PROJECT                               ID                                    TYPE    NAME  DESCRIPTION      LABELS                                                                       
            1.1.1.1  c40ad996-e1fd-4511-a7bf-418219cb8d95  2e0144a2-09ef-42b7-b629-4263295db6e8  static  ip1   ip1 description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
`),
			Template:     new("{{ .ip }} {{ .project }}"),
			WantTemplate: new(`1.1.1.1 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | IP      | PROJECT                              | ID                                   | TYPE   | NAME | ATTACHED SERVICE |
            |---------|--------------------------------------|--------------------------------------|--------|------|------------------|
            | 1.1.1.1 | c40ad996-e1fd-4511-a7bf-418219cb8d95 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static | ip1  | ingress-nginx    |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
