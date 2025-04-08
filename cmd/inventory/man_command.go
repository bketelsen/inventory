package main

import (
	"fmt"
	"os"

	"github.com/bketelsen/toolbox/cobra"
	mcoral "github.com/bketelsen/toolbox/mcobra"
	"github.com/muesli/roff"
	"github.com/spf13/viper"
)

// Build the cobra command that handles our command line tool.
func NewManCommand(config *viper.Viper) *cobra.Command {

	manCmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates inventory's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			manPage, err := mcoral.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return manCmd
}
