package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/metal-stack-cloud/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/spf13/cobra"
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
		Long:    "You can switch back and forth contexts with \"-\"",
		Example: `
~/.metal-stack-cloud/config.yaml
---
current: dev
previous: prod
contexts:
  - name: dev
	token: <dev-token>
	default-project: dev-project
  - name: prod
	token: <prod-token>
	default-project: prod-project
...
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 && args[0] == "-" {
				return w.set(args)
			}

			return cmd.Usage()
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
	contextSetCmd := &cobra.Command{
		Use:   "set",
		Short: "set the cli context",
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
		Use:   "set-project",
		Short: "sets the default project to act on for cli commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.setProject(args)
		},
		ValidArgsFunction: c.Completion.ProjectListCompletion,
	}

	contextCmd.AddCommand(
		contextListCmd,
		contextSetCmd,
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
			fmt.Printf("%s context \"%s\" already active\n", color.GreenString("✔"), color.GreenString(ctxs.CurrentContext))
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
