package v1

import (
	"fmt"
	"strings"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type volume struct {
	c *config.Config
}

func newVolumeCmd(c *config.Config) *cobra.Command {
	w := &volume{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[any, *apiv1.VolumeServiceUpdateRequest, *apiv1.Volume]{
		BinaryName:  config.BinaryName,
		GenericCLI:  genericcli.NewGenericCLI[any, *apiv1.VolumeServiceUpdateRequest, *apiv1.Volume](w).WithFS(c.Fs),
		Singular:    "volume",
		Plural:      "volumes",
		Description: "volume related actions of metalstack.cloud",
		// Sorter:          sorters.TenantSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("uuid", "", "", "filter by uuid")
			cmd.Flags().StringP("name", "", "", "filter by name")
			cmd.Flags().StringP("partition", "", "", "filter by partition")
			cmd.Flags().StringP("project", "p", "", "filter by project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
			genericcli.Must(cmd.RegisterFlagCompletionFunc("partition", c.Completion.PartitionAssetListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "filter by project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "filter by project")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringSlice("labels", nil, "the volume labels in the form of <key>=<value>")
		},
		UpdateRequestFromCLI: w.updateFromCLI,
	}

	manifestCmd := &cobra.Command{
		Use:   "manifest",
		Short: "volume manifest",
		Long:  "show detailed info about the storage cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.volumeManifest(args)
		},
	}
	manifestCmd.Flags().StringP("name", "", "restored-pv", "name of the PersistentVolume")
	manifestCmd.Flags().StringP("namespace", "", "default", "namespace for the PersistentVolume")
	manifestCmd.Flags().StringP("project", "p", "", "project")

	genericcli.Must(manifestCmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))

	encryptionSecretCmd := &cobra.Command{
		Use:   "encryptionsecret",
		Short: "volume encryptionsecret template",
		Long:  "generate volume encryptionsecret template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.volumeEncryptionSecretManifest()
		},
	}
	encryptionSecretCmd.Flags().StringP("passphrase", "", "", "passphrase")
	encryptionSecretCmd.Flags().StringP("namespace", "", "default", "namespace for the EncryptionSecret")

	return genericcli.NewCmds(cmdsConfig, manifestCmd, encryptionSecretCmd)
}

func (c *volume) Create(rq any) (*apiv1.Volume, error) {
	panic("unimplemented")
}

func (c *volume) Delete(id string) (*apiv1.Volume, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.VolumeServiceDeleteRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Volume().Delete(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to delete volumes: %w", err)
	}

	return resp.Msg.Volume, nil
}

func (c *volume) Get(id string) (*apiv1.Volume, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.VolumeServiceGetRequest{
		Uuid:    id,
		Project: c.c.GetProject(),
	}

	resp, err := c.c.Client.Apiv1().Volume().Get(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}

	return resp.Msg.Volume, nil
}

func (c *volume) List() ([]*apiv1.Volume, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.VolumeServiceListRequest{
		Project: c.c.GetProject(),
	}

	if viper.IsSet("uuid") {
		req.Uuid = pointer.Pointer(viper.GetString("uuid"))
	}
	if viper.IsSet("name") {
		req.Name = pointer.Pointer(viper.GetString("name"))
	}
	if viper.IsSet("partition") {
		req.Partition = pointer.Pointer(viper.GetString("partition"))
	}

	resp, err := c.c.Client.Apiv1().Volume().List(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes: %w", err)
	}

	return resp.Msg.Volumes, nil
}

func (c *volume) Convert(r *apiv1.Volume) (string, any, *apiv1.VolumeServiceUpdateRequest, error) {
	return helpers.EncodeProject(r.Uuid, r.Project), nil, volumeResponseToUpdate(r), nil
}

func (c *volume) Update(rq *apiv1.VolumeServiceUpdateRequest) (*apiv1.Volume, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	var (
		labels    []*apiv1.VolumeLabel
		rawLabels = viper.GetStringSlice("labels")
	)
	for _, l := range rawLabels {
		labelKey, labelValue, ok := strings.Cut(l, "=")
		if !ok {
			return nil, fmt.Errorf("label %q is not in the form of <key>=<value>", l)
		}
		labels = append(labels, &apiv1.VolumeLabel{
			Key:   labelKey,
			Value: labelValue,
		})
	}

	resp, err := c.c.Client.Apiv1().Volume().Update(ctx, connect.NewRequest(&apiv1.VolumeServiceUpdateRequest{
		Uuid:    rq.Uuid,
		Project: c.c.GetProject(),
		Labels:  labels,
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Volume, nil
}

func (c *volume) volumeManifest(args []string) error {
	id, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	volume, err := c.Get(id)
	if err != nil {
		return err
	}

	name := viper.GetString("name")
	namespace := viper.GetString("namespace")

	filesystem := corev1.PersistentVolumeFilesystem
	pv := corev1.PersistentVolume{
		TypeMeta:   v1.TypeMeta{Kind: "PersistentVolume", APIVersion: "v1"},
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			VolumeMode:  &filesystem,
			// FIXME add Capacity once figured out
			StorageClassName: volume.StorageClass,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				CSI: &corev1.CSIPersistentVolumeSource{
					Driver:       "csi.lightbitslabs.com",
					FSType:       "ext4",
					ReadOnly:     false,
					VolumeHandle: volume.VolumeHandle,
				},
			},
		},
	}

	if len(volume.AttachedTo) > 0 {
		nodes := connectedHosts(volume)
		_, _ = fmt.Fprintf(c.c.Out, "# be cautios! at the time being your volume:%s is still attached to worker node:%s, you can not mount it twice\n", volume.Uuid, strings.Join(nodes, ","))
	}

	y, err := yaml.Marshal(pv)
	if err != nil {
		panic(fmt.Errorf("unable to marshal to yaml: %w", err))
	}

	_, _ = fmt.Fprintf(c.c.Out, "---\n%s", string(y))

	return nil
}

func (v *volume) volumeEncryptionSecretManifest() error {
	namespace := viper.GetString("namespace")
	passphrase := viper.GetString("passphrase")
	secret := corev1.Secret{
		TypeMeta: v1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: v1.ObjectMeta{
			Name:      "storage-encryption-key",
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"host-encryption-passphrase": passphrase,
		},
	}
	y, err := yaml.Marshal(secret)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(v.c.Out, `# Sample secret to be used in conjunction with the partition-gold-encrypted StorageClass.
# Place this secret, after modifying namespace and the actual secret value, in the same namespace as the pvc.
#
# IMPORTANT
# Remember to make a safe copy of this secret at a secure location, once lost all your data will be lost as well.`)
	_, _ = fmt.Fprintln(v.c.Out, string(y))
	return nil
}

// connectedHosts returns the worker nodes without internal prefixes and suffixes
func connectedHosts(vol *apiv1.Volume) []string {
	nodes := []string{}
	for _, n := range vol.AttachedTo {
		// nqn.2019-09.com.lightbitslabs:host:shoot--pddhz9--duros-tst9-group-0-6b7bb-2cnvs.node
		parts := strings.Split(n, ":host:")
		if len(parts) >= 1 {
			node := strings.TrimSuffix(parts[1], ".node")
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (v *volume) updateFromCLI(args []string) (*apiv1.VolumeServiceUpdateRequest, error) {
	uuid, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	ctx, cancel := v.c.NewRequestContext()
	defer cancel()

	vol, err := v.c.Client.Apiv1().Volume().Get(ctx, connect.NewRequest(&apiv1.VolumeServiceGetRequest{
		Uuid:    uuid,
		Project: v.c.GetProject(),
	}))
	if err != nil {
		return nil, err
	}

	return &apiv1.VolumeServiceUpdateRequest{
		Uuid:    vol.Msg.Volume.Uuid,
		Project: vol.Msg.Volume.Project,
		Labels:  vol.Msg.Volume.Labels,
	}, nil
}

func volumeResponseToUpdate(r *apiv1.Volume) *apiv1.VolumeServiceUpdateRequest {
	return &apiv1.VolumeServiceUpdateRequest{
		Uuid:    r.Uuid,
		Project: r.Project,
		Labels:  r.Labels,
	}
}
