package api_e2e

import (
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

func Test_VolumeCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceListResponse, []*apiv1.Volume]{
		{
			Name: "list",
			CmdArgs: []string{
				"storage",
				"volume",
				"list",
				"--project",
				testresources.Volume1().Project,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("List", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceListRequest{
								Project: testresources.Volume1().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceListResponse{
								Volumes: []*apiv1.Volume{
									testresources.Volume1(),
								},
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
		`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
		`),
			Template: new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`
bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95
		`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
		`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Describe(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceGetResponse, *apiv1.Volume]{
		{
			Name: "describe",
			CmdArgs: []string{
				"storage",
				"volume",
				"describe",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
								Uuid: testresources.Volume1().Uuid,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
								Volume: testresources.Volume1(),
							}), nil)
						},
					},
				}}),
			WantObject: testresources.Volume1(),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Delete(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceDeleteResponse, *apiv1.Volume]{
		{
			Name: "delete",
			CmdArgs: []string{
				"storage",
				"volume",
				"delete",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceDeleteRequest{
								Uuid: testresources.Volume1().Uuid,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceDeleteResponse{
								Volume: testresources.Volume1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Volume1(),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
`),
		},
		/* {
					// Msg: (*apiv1.VolumeServiceDeleteRequest)(uuid:"bd0f32e2-eabf-4eb7-a0db-25fc993c3678 (c40ad996-e1fd-4511-a7bf-418219cb8d95)")
					Name:    "delete volume from file",
					CmdArgs: append([]string{"storage", "volume", "delete"}, e2e.AppendFromFileCommonArgs()...),
					NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
						FsMocks: func(fs *afero.Afero) {
							require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Volume1()), 0755))
						},
						ClientMocks: &apitests.ClientMockFns{
							Apiv1Mocks: &apitests.Apiv1MockFns{
								Volume: func(m *mock.Mock) {
									m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
										Uuid:    testresources.Volume1().Uuid,
										Project: testresources.Volume1().Project,
									})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
										Volume: testresources.Volume1(),
									}), nil)
									m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceDeleteRequest{
										Uuid:    testresources.Volume1().Uuid,
										Project: testresources.Volume1().Project,
									})).Return(connect.NewResponse(&apiv1.VolumeServiceDeleteResponse{
										Volume: testresources.Volume1(),
									}), nil)
								},
							},
						},
					}),
					WantTable: new(`
				            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION
				            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
									`),
				},
				{
					// Msg: (*apiv1.VolumeServiceDeleteRequest)(uuid:"bd0f32e2-eabf-4eb7-a0db-25fc993c3678 (c40ad996-e1fd-4511-a7bf-418219cb8d95)")
					Name:    "delete many volume from file",
					CmdArgs: append([]string{"storage", "volume", "delete"}, e2e.AppendFromFileCommonArgs()...),
					NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
						FsMocks: func(fs *afero.Afero) {
							require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Volume1(), testresources.Volume2()), 0755))
						},
						ClientMocks: &apitests.ClientMockFns{
							Apiv1Mocks: &apitests.Apiv1MockFns{
								Volume: func(m *mock.Mock) {
									m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
										Uuid:    testresources.Volume2().Uuid,
										Project: testresources.Volume2().Project,
									})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
										Volume: testresources.Volume2(),
									}), nil)
									m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceDeleteRequest{
										Uuid: testresources.Volume2().Uuid,
									})).Return(connect.NewResponse(&apiv1.VolumeServiceDeleteResponse{
										Volume: testresources.Volume2(),
									}), nil)
									m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
										Uuid:    testresources.Volume1().Uuid,
										Project: testresources.Volume1().Project,
									})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
										Volume: testresources.Volume1(),
									}), nil)
									m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceDeleteRequest{
										Uuid: testresources.Volume1().Uuid,
									})).Return(connect.NewResponse(&apiv1.VolumeServiceDeleteResponse{
										Volume: testresources.Volume1(),
									}), nil)
								},
							},
						},
					}),
					WantTable: new(`
		            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION
		            0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster2      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d67  partition-a
		            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
							`),
				}, */
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Update(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceUpdateResponse, *apiv1.Volume]{
		{
			Name: "update volume",
			CmdArgs: []string{
				"storage",
				"volume",
				"update",
				testresources.Volume1().Uuid,
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
								Uuid: testresources.Volume1().Uuid,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceUpdateResponse{
								Volume: testresources.Volume1(),
							}), nil)
						},
					},
				},
			}),
			WantObject: testresources.Volume1(),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
`),
			WantWideTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    NODES  LABELS   
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a         foo=bar
`),
			Template:     new("{{ .uuid }} {{ .project }}"),
			WantTemplate: new(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 c40ad996-e1fd-4511-a7bf-418219cb8d95`),
			WantMarkdown: new(`
            | ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT                              | PARTITION   |
            |--------------------------------------|---------|---------|-------|----------|--------------|----------------|--------------------------------------|-------------|
            | bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster1     | storageclass-a | c40ad996-e1fd-4511-a7bf-418219cb8d95 | partition-a |
`),
		},
		{
			Name:    "update volume from file",
			CmdArgs: append([]string{"storage", "volume", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshal(t, testresources.Volume1()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
								Uuid:    testresources.Volume1().Uuid,
								Project: testresources.Volume1().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
								Volume: testresources.Volume1(),
							}), nil)
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
								Uuid:    testresources.Volume1().Uuid,
								Project: testresources.Volume1().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceUpdateResponse{
								Volume: testresources.Volume1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
					`),
		},
		{
			Name:    "update many volume from file",
			CmdArgs: append([]string{"storage", "volume", "update"}, e2e.AppendFromFileCommonArgs()...),
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				FsMocks: func(fs *afero.Afero) {
					require.NoError(t, fs.WriteFile(e2e.InputFilePath, e2e.MustMarshalToMultiYAML(t, testresources.Volume1(), testresources.Volume2()), 0755))
				},
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
								Uuid:    testresources.Volume2().Uuid,
								Project: testresources.Volume2().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
								Volume: testresources.Volume2(),
							}), nil)
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
								Uuid:    testresources.Volume2().Uuid,
								Project: testresources.Volume2().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceUpdateResponse{
								Volume: testresources.Volume2(),
							}), nil)
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
								Uuid:    testresources.Volume1().Uuid,
								Project: testresources.Volume1().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
								Volume: testresources.Volume1(),
							}), nil)
							m.On("Update", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
								Uuid:    testresources.Volume1().Uuid,
								Project: testresources.Volume1().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceUpdateResponse{
								Volume: testresources.Volume1(),
							}), nil)
						},
					},
				},
			}),
			WantTable: new(`
            ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT                               PARTITION    
            0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster2      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d67  partition-a  
            bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster1      storageclass-a  c40ad996-e1fd-4511-a7bf-418219cb8d95  partition-a
					`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_Manifest(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceGetResponse, string]{
		{
			Name: "volume manifest",
			CmdArgs: []string{
				"storage",
				"volume",
				"manifest",
				testresources.Volume1().Uuid,
				"--project", testresources.Volume1().Project,
				"--name", "test-name",
				"--namespace", "test-namespace",
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientMocks: &apitests.ClientMockFns{
					Apiv1Mocks: &apitests.Apiv1MockFns{
						Volume: func(m *mock.Mock) {
							m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
								Uuid:    testresources.Volume1().Uuid,
								Project: testresources.Volume1().Project,
							})).Return(connect.NewResponse(&apiv1.VolumeServiceGetResponse{
								Volume: testresources.Volume1(),
							}), nil)
						},
					},
				}}),
			WantDefault: new(`
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: test-name
  namespace: test-namespace
spec:
  accessModes:
  - ReadWriteOnce
  csi:
    driver: csi.lightbitslabs.com
    fsType: ext4
    volumeHandle: ""
  storageClassName: storageclass-a
  volumeMode: Filesystem
status: {}
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_EncryptionSecret(t *testing.T) {
	tests := []*e2e.Test[apiv1.VolumeServiceGetResponse, string]{
		{
			Name: "volume manifest",
			CmdArgs: []string{
				"storage",
				"volume",
				"encryptionsecret",
				testresources.Volume1().Uuid,
				"--passphrase", "test-phrase",
				"--namespace", "test-namespace",
			},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{}),
			WantDefault: new(`
# Sample secret to be used in conjunction with the partition-gold-encrypted StorageClass.
# Place this secret, after modifying namespace and the actual secret value, in the same namespace as the pvc.
#
# IMPORTANT
# Remember to make a safe copy of this secret at a secure location, once lost all your data will be lost as well.apiVersion: v1
kind: Secret
metadata:
  name: storage-encryption-key
  namespace: test-namespace
stringData:
  host-encryption-passphrase: test-phrase
type: Opaque
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
