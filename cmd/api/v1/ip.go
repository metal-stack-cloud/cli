package v1

import (
	"fmt"

	"connectrpc.com/connect"
	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack-cloud/cli/cmd/sorters"
	"github.com/metal-stack-cloud/cli/pkg/helpers"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/genericcli/printers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ip struct {
	c *config.Config
}

func newIPCmd(c *config.Config) *cobra.Command {
	w := &ip{
		c: c,
	}

	cmdsConfig := &genericcli.CmdsConfig[*apiv1.IPServiceAllocateRequest, *apiv1.IPServiceUpdateRequest, *apiv1.IP]{
		BinaryName:      config.BinaryName,
		GenericCLI:      genericcli.NewGenericCLI[*apiv1.IPServiceAllocateRequest, *apiv1.IPServiceUpdateRequest, *apiv1.IP](w).WithFS(c.Fs),
		Singular:        "ip",
		Plural:          "ips",
		Description:     "an ip address of metalstack.cloud",
		Sorter:          sorters.IPSorter(),
		DescribePrinter: func() printers.Printer { return c.DescribePrinter },
		ListPrinter:     func() printers.Printer { return c.ListPrinter },
		ListCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project from where ips should be listed")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		CreateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")
			cmd.Flags().StringP("name", "", "", "name of the ip")
			cmd.Flags().StringP("description", "", "", "description of the ip")
			cmd.Flags().StringSliceP("tags", "", nil, "tags to add to the ip")
			cmd.Flags().BoolP("static", "", false, "make this ip static")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		UpdateCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")
			cmd.Flags().String("name", "", "name of the ip")
			cmd.Flags().String("description", "", "description of the ip")
			cmd.Flags().StringSlice("tags", nil, "tags of the ip")
			cmd.Flags().Bool("static", false, "make this ip static")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DescribeCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		DeleteCmdMutateFn: func(cmd *cobra.Command) {
			cmd.Flags().StringP("project", "p", "", "project of the ip")

			genericcli.Must(cmd.RegisterFlagCompletionFunc("project", c.Completion.ProjectListCompletion))
		},
		CreateRequestFromCLI: func() (*apiv1.IPServiceAllocateRequest, error) {
			return &apiv1.IPServiceAllocateRequest{
				Project:     c.GetProject(),
				Name:        viper.GetString("name"),
				Description: viper.GetString("description"),
				Tags:        viper.GetStringSlice("tags"),
				Static:      viper.GetBool("static"),
			}, nil
		},
		UpdateRequestFromCLI: w.updateFromCLI,
		ValidArgsFn:          c.Completion.IpListCompletion,
	}

	return genericcli.NewCmds(cmdsConfig)
}

func (c *ip) updateFromCLI(args []string) (*apiv1.IPServiceUpdateRequest, error) {
	uuid, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return nil, err
	}

	ipToUpdate, err := c.Get(uuid)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve ip: %w", err)
	}

	if viper.IsSet("name") {
		ipToUpdate.Name = viper.GetString("name")
	}
	if viper.IsSet("description") {
		ipToUpdate.Description = viper.GetString("description")
	}
	if viper.IsSet("static") {
		ipToUpdate.Type = ipStaticToType(viper.GetBool("static"))
	}
	if viper.IsSet("tags") {
		ipToUpdate.Tags = viper.GetStringSlice("tags")
	}

	return IpResponseToUpdate(ipToUpdate), nil
}

func (c *ip) Create(rq *apiv1.IPServiceAllocateRequest) (*apiv1.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().IP().Allocate(ctx, connect.NewRequest(rq))
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			return nil, genericcli.AlreadyExistsError()
		}
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (c *ip) Delete(id string) (*apiv1.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &apiv1.IPServiceDeleteRequest{
		Project: c.c.GetProject(),
		Uuid:    id,
	}

	if viper.IsSet("file") {
		var err error
		req.Uuid, req.Project, err = helpers.DecodeProject(id)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Client.Apiv1().IP().Delete(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (c *ip) Get(id string) (*apiv1.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().IP().Get(ctx, connect.NewRequest(&apiv1.IPServiceGetRequest{
		Project: c.c.GetProject(),
		Uuid:    id,
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (c *ip) List() ([]*apiv1.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().IP().List(ctx, connect.NewRequest(&apiv1.IPServiceListRequest{
		Project: c.c.GetProject(),
	}))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ips, nil
}

func (c *ip) Update(rq *apiv1.IPServiceUpdateRequest) (*apiv1.IP, error) {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	resp, err := c.c.Client.Apiv1().IP().Update(ctx, connect.NewRequest(rq))
	if err != nil {
		return nil, err
	}

	return resp.Msg.Ip, nil
}

func (*ip) Convert(r *apiv1.IP) (string, *apiv1.IPServiceAllocateRequest, *apiv1.IPServiceUpdateRequest, error) {
	return helpers.EncodeProject(r.Uuid, r.Project), IpResponseToCreate(r), IpResponseToUpdate(r), nil
}

func IpResponseToCreate(r *apiv1.IP) *apiv1.IPServiceAllocateRequest {
	return &apiv1.IPServiceAllocateRequest{
		Project:     r.Project,
		Name:        r.Name,
		Description: r.Description,
		Tags:        r.Tags,
		Static:      ipTypeToStatic(r.Type),
	}
}

func IpResponseToUpdate(r *apiv1.IP) *apiv1.IPServiceUpdateRequest {
	return &apiv1.IPServiceUpdateRequest{
		Project: r.Project,
		Ip:      r,
	}
}

func ipStaticToType(b bool) apiv1.IPType {
	if b {
		return apiv1.IPType_IP_TYPE_STATIC
	}
	return apiv1.IPType_IP_TYPE_EPHEMERAL
}

func ipTypeToStatic(t apiv1.IPType) bool {
	return t == apiv1.IPType_IP_TYPE_STATIC
}
