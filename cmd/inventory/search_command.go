package main

import (
	"fmt"

	"github.com/bketelsen/inventory/client"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/viper"
)

// Build the cobra command that handles our command line tool.
func NewSearchCommand(config *viper.Viper) *cobra.Command {

	// Define our command
	searchCmd := &cobra.Command{
		Use:          "search [query]",
		SilenceUsage: true,
		Example: `inventory search jellyfin
	// more verbose output
	inventory search --verbose jellyfin
	// specify a config file
	inventory search --verbose --config /path/to/config.yaml jellyfin`,
		Short: "search for services or containers",
		Args:  cobra.ExactArgs(1),
		Long:  `Search returns a list of services, listeners, and containers that match the query.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Get server address from config or use default
			cl := client.NewClient(config.GetString("client.remote"), config.GetString("client.location"), config.GetString("client.description"), nil)

			result, err := cl.Search(args[0])
			if err != nil {
				ui.Error("Error searching", err.Error())
				return err
			}
			if len(result) == 0 {
				ui.Info("No results found")
				return nil
			}
			ui.Info(fmt.Sprintf("Found %d Server(s) with \"%s\"", len(result), args[0]))

			for _, r := range result {
				cmd.Printf("Server %s / %s:\n", r.Host.HostName, r.Host.IP)
				if len(r.Listeners) > 0 {
					cmd.Printf("Found %d Listener(s) on %s / %s:\n", len(r.Listeners), r.Host.HostName, r.Host.IP)

					out, err := ui.DisplayTable(r.Listeners, "", nil)
					if err != nil {
						ui.Error("Error displaying results", err.Error())
					}
					cmd.Println(out)
					cmd.Println()
				}
				if len(r.Services) > 0 {
					cmd.Printf("Found %d Service(s) on %s / %s:\n", len(r.Services), r.Host.HostName, r.Host.IP)
					out, err := ui.DisplayTable(r.Services, "", nil)
					if err != nil {
						ui.Error("Error displaying results", err.Error())
					}
					cmd.Println(out)
					cmd.Println()
				}

				if len(r.Containers) > 0 {
					cmd.Printf("Found %d Containers on %s / %s:\n", len(r.Containers), r.Host.HostName, r.Host.IP)

					out, err := ui.DisplayTable(r.Containers, "", nil)
					if err != nil {
						ui.Error("Error displaying results", err.Error())
					}
					cmd.Println(out)
					cmd.Println()
				}
			}
			// out, err := ui.DisplayTable(result, "", nil)
			// if err != nil {
			// 	cmd.Println(ui.Error("Error displaying results", err.Error()))
			// }
			// cmd.Println(out)
			return nil

		},
	}
	searchCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		_ = config.BindPFlag("client.remote", cmd.Flags().Lookup("remote"))
		return nil

	}
	// Define cobra flags, the default value has the lowest (least significant) precedence
	searchCmd.Flags().StringP("remote", "r", "10.0.1.1:9999", "Remote inventory server address")

	return searchCmd
}
