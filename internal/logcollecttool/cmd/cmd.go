package cmd

import (
	"flag"
	"io"
	"logconvert/internal/logcollecttool/cmd/options"
	"logconvert/internal/logcollecttool/cmd/strategy"
	"logconvert/internal/logcollecttool/util/templates"
	"os"

	"github.com/spf13/cobra"
	cliflag "logconvert/cli/flag"
	"logconvert/cli/genericclioptions"
	rootoptions "logconvert/internal/pkg/options"
)

func NewDefaultLogCollectToolCommand() *cobra.Command {
	return NewLogCollectToolCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewLogCollectToolCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:               "logcollecttool",
		Short:             "logcollecttool locate and analyze problems for logcollect",
		Long:              templates.LongDesc(`logcollecttool locate and analyze problems for logcollect.`),
		Run:               runHelp,
		DisableAutoGenTag: true,
	}

	flags := cmds.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc)

	rootOptions := rootoptions.NewServerRunOptions()
	rootOptions.AddFlags(flags)

	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	cmds.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	groups := templates.CommandGroups{
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				strategy.NewCmdValidate(rootOptions, ioStreams),
			},
		},
	}
	groups.Add(cmds)

	filters := []string{"options", "completion"}
	templates.ActsAsRootCommand(cmds, filters, groups...)

	cmds.AddCommand(options.NewCmdOptions(ioStreams.Out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
