package cmd

import (
	"fmt"
	"os"

	"github.com/bketelsen/toolbox/cobra"
	mcoral "github.com/bketelsen/toolbox/mcobra"
	"github.com/muesli/roff"
)

type manCmd struct {
	cmd *cobra.Command
}

func newManCmd() *manCmd {
	root := &manCmd{}
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates inventory's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			manPage, err := mcoral.NewManPage(1, root.cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	root.cmd = cmd
	return root
}

func init() {
	manCmd := newManCmd()
	rootCmd.AddCommand(manCmd.cmd)
}
