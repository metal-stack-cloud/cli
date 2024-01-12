package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ctx struct {
	c *config.Config
}

func newContextCmd(c *config.Config) *cobra.Command {
	w := &ctx{
		c: c,
	}

	contextCmd := &cobra.Command{
		Use:     "context",
		Aliases: []string{"ctx"},
		Short:   "manage cli contexts",
		Long:    "you can switch back and forth contexts with \"-\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.set(args)
		},
	}

	contextListCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list the configured cli contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.list()
		},
	}
	contextSwitchCmd := &cobra.Command{
		Use:     "switch <context-name>",
		Short:   "switch the cli context",
		Long:    "you can switch back and forth contexts with \"-\"",
		Aliases: []string{"set"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.set(args)
		},
		ValidArgsFunction: config.ContextListCompletion,
	}
	contextShortCmd := &cobra.Command{
		Use:   "show-current",
		Short: "prints the current context name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.short()
		},
	}
	contextSetProjectCmd := &cobra.Command{
		Use:   "set-project <project-id>",
		Short: "sets the default project to act on for cli commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.setProject(args)
		},
		ValidArgsFunction: c.Completion.ProjectListCompletion,
	}
	contextRemoveCmd := &cobra.Command{
		Use:     "remove <context-name>",
		Aliases: []string{"rm", "delete"},
		Short:   "remove a cli context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.remove(args)
		},
		ValidArgsFunction: config.ContextListCompletion,
	}

	contextAddCmd := &cobra.Command{
		Use:     "add <context-name>",
		Aliases: []string{"create"},
		Short:   "add a cli context",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.add(args)
		},
	}
	contextAddCmd.Flags().String("api-url", "", "sets the api-url for this context")
	contextAddCmd.Flags().String("api-token", "", "sets the api-token for this context")
	contextAddCmd.Flags().String("default-project", "", "sets a default project to act on")
	contextAddCmd.Flags().Duration("timeout", 0, "sets a default request timeout")
	contextAddCmd.Flags().Bool("activate", false, "immediately switches to the new context")

	genericcli.Must(contextAddCmd.MarkFlagRequired("api-token"))

	contextCmd.AddCommand(
		contextListCmd,
		contextSwitchCmd,
		contextAddCmd,
		contextRemoveCmd,
		contextShortCmd,
		contextSetProjectCmd,
	)

	return contextCmd
}

func (c *ctx) list() error {
	ctxs, err := config.GetContexts()
	if err != nil {
		return err
	}

	return c.c.ListPrinter.Print(ctxs)
}

func (c *ctx) short() error {
	ctxs, err := config.GetContexts()
	if err != nil {
		return err
	}
	fmt.Println(ctxs.CurrentContext)
	return nil
}

func (c *ctx) add(args []string) error {
	name, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return fmt.Errorf("no context name given")
	}

	ctxs, err := config.GetContexts()
	if err != nil {
		return err
	}

	ctx := &config.Context{
		Name:           name,
		ApiURL:         pointer.PointerOrNil(viper.GetString("api-url")),
		Token:          viper.GetString("api-token"),
		DefaultProject: viper.GetString("default-project"),
		Timeout:        pointer.PointerOrNil(viper.GetDuration("timeout")),
	}

	ctxs.Contexts = append(ctxs.Contexts, ctx)

	if viper.GetBool("activate") {
		ctxs.PreviousContext = ctxs.CurrentContext
		ctxs.CurrentContext = ctx.Name
	}

	err = config.WriteContexts(ctxs)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.c.Out, "%s added context \"%s\"\n", color.GreenString("✔"), color.GreenString(ctx.Name))

	return nil
}

func (c *ctx) remove(args []string) error {
	name, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return fmt.Errorf("no context name given")
	}

	ctxs, err := config.GetContexts()
	if err != nil {
		return err
	}

	ctx, ok := ctxs.GetContext(name)
	if !ok {
		return fmt.Errorf("no context with name %q found", name)
	}

	ctxs.Delete(ctx.Name)

	err = config.WriteContexts(ctxs)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.c.Out, "%s removed context \"%s\"\n", color.GreenString("✔"), color.GreenString(ctx.Name))

	return nil
}

func (c *ctx) set(args []string) error {
	wantCtx, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return fmt.Errorf("no context name given")
	}

	ctxs, err := config.GetContexts()
	if err != nil {
		return err
	}

	if wantCtx == "-" {
		prev := ctxs.PreviousContext
		if prev == "" {
			return fmt.Errorf("no previous context found")
		}

		curr := ctxs.CurrentContext
		ctxs.PreviousContext = curr
		ctxs.CurrentContext = prev
	} else {
		nextCtx := wantCtx
		_, ok := ctxs.GetContext(nextCtx)
		if !ok {
			return fmt.Errorf("context %s not found", nextCtx)
		}
		if nextCtx == ctxs.CurrentContext {
			fmt.Fprintf(c.c.Out, "%s context \"%s\" already active\n", color.GreenString("✔"), color.GreenString(ctxs.CurrentContext))
			return nil
		}
		ctxs.PreviousContext = ctxs.CurrentContext
		ctxs.CurrentContext = nextCtx
	}

	err = config.WriteContexts(ctxs)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.c.Out, "%s switched context to \"%s\"\n", color.GreenString("✔"), color.GreenString(ctxs.CurrentContext))

	return nil
}

func (c *ctx) setProject(args []string) error {
	project, err := genericcli.GetExactlyOneArg(args)
	if err != nil {
		return err
	}

	ctxs, err := config.GetContexts()
	if err != nil {
		return err
	}

	ctx, ok := ctxs.GetContext(c.c.Context.Name)
	if !ok {
		return fmt.Errorf("no context currently active")
	}

	ctx.DefaultProject = project

	err = config.WriteContexts(ctxs)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.c.Out, "%s switched context default project to \"%s\"\n", color.GreenString("✔"), color.GreenString(ctx.DefaultProject))

	return nil
}
