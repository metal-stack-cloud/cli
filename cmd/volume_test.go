package cmd

import (
	"testing"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	apitests "github.com/metal-stack-cloud/api/go/tests"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/stretchr/testify/mock"
)

var (
	volume1 = func() *apiv1.Volume {
		return volume1WithLabels([]*apiv1.VolumeLabel{
			{
				Key:   "foo",
				Value: "bar",
			},
		})
	}
	volume1WithLabels = func(labels []*apiv1.VolumeLabel) *apiv1.Volume {
		return &apiv1.Volume{
			Uuid:               "bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
			Name:               "volume1",
			Project:            "a",
			Partition:          "partition-a",
			StorageClass:       "storageclass-a",
			Size:               1024,
			Usage:              42,
			State:              "Bound",
			AttachedTo:         nil,
			SourceSnapshotUuid: "",
			SourceSnapshotName: "",
			VolumeHandle:       "",
			NodeIps:            nil,
			RebuildProgress:    "",
			PrimaryNodeUuid:    "",
			QosPolicyUuid:      "",
			QosPolicyName:      "",
			ReplicaCount:       0,
			ProtectionState:    "",
			LogicalUsedStorage: 0,
			Statistics:         nil,
			ClusterName:        "cluster-a1",
			ClusterId:          "",
			Labels:             labels,
		}
	}
	volume2 = func() *apiv1.Volume {
		return &apiv1.Volume{
			Uuid:               "0372d029-1077-4e9b-b303-7d64ad5496fd",
			Name:               "volume2",
			Project:            "a",
			Partition:          "partition-a",
			StorageClass:       "storageclass-a",
			Size:               1024,
			Usage:              42,
			State:              "Bound",
			AttachedTo:         nil,
			SourceSnapshotUuid: "",
			SourceSnapshotName: "",
			VolumeHandle:       "",
			NodeIps:            nil,
			RebuildProgress:    "",
			PrimaryNodeUuid:    "",
			QosPolicyUuid:      "",
			QosPolicyName:      "",
			ReplicaCount:       0,
			ProtectionState:    "",
			LogicalUsedStorage: 0,
			Statistics:         nil,
			ClusterName:        "cluster-a2",
			ClusterId:          "",
			Labels: []*apiv1.VolumeLabel{
				{
					Key:   "bar",
					Value: "baz",
				},
			},
		}
	}
)

func Test_VolumeCmd_MultiResult(t *testing.T) {
	tests := []*Test[[]*apiv1.Volume]{
		{
			Name: "list",
			Cmd: func(want []*apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"list",
					"--project",
					"a",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("List", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceListRequest{
							Project: "a",
						})).Return(&connect.Response[apiv1.VolumeServiceListResponse]{Msg: &apiv1.VolumeServiceListResponse{
							Volumes: []*apiv1.Volume{
								volume2(),
								volume1(),
							},
						},
						}, nil)
					},
				},
			},
			Want: []*apiv1.Volume{
				volume1(),
				volume2(),
			},
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION
0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster-a2    storageclass-a  a        partition-a
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a
		`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS
0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster-a2    storageclass-a  a        partition-a         bar=baz
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a         foo=bar
		`),
			Template: pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
bd0f32e2-eabf-4eb7-a0db-25fc993c3678 a
0372d029-1077-4e9b-b303-7d64ad5496fd a
		`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| 0372d029-1077-4e9b-b303-7d64ad5496fd | volume2 | 1.0 KiB | 42 B  | 0        | cluster-a2   | storageclass-a | a       | partition-a |
| bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster-a1   | storageclass-a | a       | partition-a |
		`),
		},
		{
			Name: "list reverse order",
			Cmd: func(want []*apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"list",
					"--project",
					"a",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("List", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceListRequest{
							Project: "a",
						})).Return(&connect.Response[apiv1.VolumeServiceListResponse]{Msg: &apiv1.VolumeServiceListResponse{
							Volumes: []*apiv1.Volume{
								volume1(),
								volume2(),
							},
						},
						}, nil)
					},
				},
			},
			Want: []*apiv1.Volume{
				volume1(),
				volume2(),
			},
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    
0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster-a2    storageclass-a  a        partition-a  
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a
		`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS   
0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster-a2    storageclass-a  a        partition-a         bar=baz  
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a         foo=bar
		`),
			Template: pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
bd0f32e2-eabf-4eb7-a0db-25fc993c3678 a
0372d029-1077-4e9b-b303-7d64ad5496fd a
		`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| 0372d029-1077-4e9b-b303-7d64ad5496fd | volume2 | 1.0 KiB | 42 B  | 0        | cluster-a2   | storageclass-a | a       | partition-a |
| bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster-a1   | storageclass-a | a       | partition-a |
		`),
		},

		{
			Name: "list just volume1",
			Cmd: func(want []*apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"list",
					"--project",
					"a",
					"--name",
					"volume1",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("List", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceListRequest{
							Project: "a",
							Name:    pointer.Pointer("volume1"),
						})).Return(&connect.Response[apiv1.VolumeServiceListResponse]{Msg: &apiv1.VolumeServiceListResponse{
							Volumes: []*apiv1.Volume{
								volume1(),
							},
						},
						}, nil)
					},
				},
			},
			Want: []*apiv1.Volume{
				volume1(),
			},
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a
		`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a         foo=bar
		`),
			Template: pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`
bd0f32e2-eabf-4eb7-a0db-25fc993c3678 a
		`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster-a1   | storageclass-a | a       | partition-a |
		`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_VolumeCmd_SingleResult(t *testing.T) {
	tests := []*Test[*apiv1.Volume]{
		{
			Name: "describe volume1",
			Cmd: func(want *apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"describe",
					"bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
							Uuid: "bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
						})).Return(&connect.Response[apiv1.VolumeServiceGetResponse]{Msg: &apiv1.VolumeServiceGetResponse{
							Volume: volume1(),
						},
						}, nil)
					},
				},
			},
			Want: volume1(),
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a
`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS   
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a         foo=bar
`),
			Template:     pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 a`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster-a1   | storageclass-a | a       | partition-a |
`),
		},
		{
			Name: "describe volume2",
			Cmd: func(want *apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"describe",
					"0372d029-1077-4e9b-b303-7d64ad5496fd",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("Get", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
							Uuid: "0372d029-1077-4e9b-b303-7d64ad5496fd",
						})).Return(&connect.Response[apiv1.VolumeServiceGetResponse]{Msg: &apiv1.VolumeServiceGetResponse{
							Volume: volume2(),
						},
						}, nil)
					},
				},
			},
			Want: volume2(),
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    
0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster-a2    storageclass-a  a        partition-a
`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS   
0372d029-1077-4e9b-b303-7d64ad5496fd  volume2  1.0 KiB  42 B   0         cluster-a2    storageclass-a  a        partition-a         bar=baz
`),
			Template:     pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`0372d029-1077-4e9b-b303-7d64ad5496fd a`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| 0372d029-1077-4e9b-b303-7d64ad5496fd | volume2 | 1.0 KiB | 42 B  | 0        | cluster-a2   | storageclass-a | a       | partition-a |
`),
		},
		{
			Name: "delete volume1",
			Cmd: func(want *apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"delete",
					"bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("Delete", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceDeleteRequest{
							Uuid: "bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
						})).Return(&connect.Response[apiv1.VolumeServiceDeleteResponse]{Msg: &apiv1.VolumeServiceDeleteResponse{
							Volume: volume1(),
						},
						}, nil)
					},
				},
			},
			Want: volume1(),
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a
`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS   
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a         foo=bar
`),
			Template:     pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 a`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster-a1   | storageclass-a | a       | partition-a |
`),
		},
		{
			Name: "update volume1",
			Cmd: func(want *apiv1.Volume) []string {
				return []string{
					"storage",
					"volume",
					"update",
					"bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
					"--add-label",
					"hello=world",
					"--remove-label",
					"foo",
				}
			},
			ClientMocks: &apitests.ClientMockFns{
				Apiv1Mocks: &apitests.Apiv1MockFns{
					Volume: func(m *mock.Mock) {
						m.On("Update", mock.Anything, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
							Uuid: "bd0f32e2-eabf-4eb7-a0db-25fc993c3678",
							Labels: &apiv1.UpdateVolumeLabels{
								Update: []*apiv1.VolumeLabel{
									{Key: "hello", Value: "world"},
								},
								Remove: []string{"foo"},
							},
						})).Return(&connect.Response[apiv1.VolumeServiceUpdateResponse]{Msg: &apiv1.VolumeServiceUpdateResponse{
							Volume: volume1WithLabels([]*apiv1.VolumeLabel{
								{
									Key:   "hello",
									Value: "world",
								},
							}),
						},
						}, nil)
					},
				},
			},
			Want: volume1WithLabels([]*apiv1.VolumeLabel{
				{
					Key:   "hello",
					Value: "world",
				},
			}),
			WantTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a
`),
			WantWideTable: pointer.Pointer(`
ID                                    NAME     SIZE     USAGE  REPLICAS  CLUSTER NAME  STORAGE CLASS   PROJECT  PARTITION    NODES  LABELS       
bd0f32e2-eabf-4eb7-a0db-25fc993c3678  volume1  1.0 KiB  42 B   0         cluster-a1    storageclass-a  a        partition-a         hello=world
`),
			Template:     pointer.Pointer("{{ .uuid }} {{ .project }}"),
			WantTemplate: pointer.Pointer(`bd0f32e2-eabf-4eb7-a0db-25fc993c3678 a`),
			WantMarkdown: pointer.Pointer(`
| ID                                   | NAME    | SIZE    | USAGE | REPLICAS | CLUSTER NAME | STORAGE CLASS  | PROJECT | PARTITION   |
|--------------------------------------|---------|---------|-------|----------|--------------|----------------|---------|-------------|
| bd0f32e2-eabf-4eb7-a0db-25fc993c3678 | volume1 | 1.0 KiB | 42 B  | 0        | cluster-a1   | storageclass-a | a       | partition-a |
`),
		},
	}

	for _, tt := range tests {
		tt.TestCmd(t)
	}
}
