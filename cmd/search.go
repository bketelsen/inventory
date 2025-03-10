/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/bketelsen/inventory/client"
	"github.com/bketelsen/inventory/types"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use: "search [query]",
	Example: `inventory search jellyfin
// more verbose output
inventory search --verbose jellyfin
// specify a config file
inventory search --verbose --config /path/to/config.yaml jellyfin`,
	Short: "search for services or containers",
	Args:  cobra.ExactArgs(1),
	Long:  `Search returns a list of services, listeners, and containers that match the query.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.SetDefault(cmd.Logger)
		// Read config
		config, err := types.ReadConfig()
		if err != nil {
			log.Println("Error reading config:", err)
			return err
		}

		// Get server address from config or use default

		cl := client.NewClient(config)
		result, err := cl.Search(args[0])
		if err != nil {
			cmd.Println(ui.Error("Error searching", err.Error()))
			return err
		}
		if len(result) == 0 {
			cmd.Println(ui.Info("No results found"))
			return nil
		}
		cmd.Println(ui.Info(fmt.Sprintf("Found %d server with \"%s\"", len(result), args[0])))

		for _, r := range result {
			if len(r.Listeners) > 0 {
				cmd.Printf("Found %d Listeners on %s / %s:\n", len(r.Listeners), r.Host.HostName, r.Host.IP)

				out, err := ui.DisplayTable(r.Listeners, "", nil)
				if err != nil {
					cmd.Println(ui.Error("Error displaying results", err.Error()))
				}
				cmd.Println(out)
				cmd.Println()
			}
			if len(r.Services) > 0 {
				cmd.Printf("Found %d Services on %s / %s:\n", len(r.Services), r.Host.HostName, r.Host.IP)
				out, err := ui.DisplayTable(r.Services, "", nil)
				if err != nil {
					cmd.Println(ui.Error("Error displaying results", err.Error()))
				}
				cmd.Println(out)
				cmd.Println()
			}

			if len(r.Containers) > 0 {
				cmd.Printf("Found %d Containers on %s / %s:\n", len(r.Containers), r.Host.HostName, r.Host.IP)

				out, err := ui.DisplayTable(r.Containers, "", nil)
				if err != nil {
					cmd.Println(ui.Error("Error displaying results", err.Error()))
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

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
