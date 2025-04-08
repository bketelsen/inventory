package main

import (
	"log"

	"github.com/bketelsen/inventory"
	"github.com/bketelsen/inventory/client"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/spf13/viper"
)

// Build the cobra command that handles our command line tool.
func NewSendCommand(config *viper.Viper) *cobra.Command {

	// Define our command
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send host and container information to the server",
		Example: `inventory send
	// more verbose output
	inventory send --verbose
	// specify a config file
	inventory send --verbose --config /path/to/config.yaml`,

		Long: `Send host and container information to the server
	This command collects information about the host and docker/incus containers
	and sends it to the server. It is designed to be run as a cron job or systemd timer.
	It is not intended to be run interactively.
	`,

		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.Logger.Info("Sending inventory")
			cmd.Logger.Info("send", "remote", config.GetString("client.remote"), "location", config.GetString("client.location"), "description", config.GetString("client.description"))

			var services []*inventory.Service
			if config.IsSet("services") {
				err := config.UnmarshalKey("services", &services)
				if err != nil {
					cmd.Logger.Error("Error getting services from config", "error", err)
					return err
				}
			}
			cl := client.NewClient(config.GetString("client.remote"), config.GetString("client.location"), config.GetString("client.description"), services)
			err := cl.Send()
			if err != nil {
				log.Printf("Error sending report: %v", err)
				return err
			}
			cmd.Logger.Info("Inventory sent successfully")

			return nil
		},
	}
	sendCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		_ = config.BindPFlag("client.remote", cmd.Flags().Lookup("remote"))
		_ = config.BindPFlag("client.location", cmd.Flags().Lookup("location"))
		_ = config.BindPFlag("client.description", cmd.Flags().Lookup("description"))
		return nil

	}
	// Define cobra flags, the default value has the lowest (least significant) precedence
	sendCmd.Flags().StringP("remote", "r", "10.0.1.1:9999", "Remote inventory server address")
	sendCmd.Flags().StringP("location", "l", "My Home Lab", "Location of the server")
	sendCmd.Flags().StringP("description", "d", "My Generic Server", "Description of the server")

	return sendCmd
}
