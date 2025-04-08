package main

import (
	"os"

	"github.com/bketelsen/inventory"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/viper"
)

// NewConfigCommand creates a new command to generate an example configuration file
func NewConfigCommand(config *viper.Viper) *cobra.Command {

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Create an configuration file for the inventory application",
		Long: `Create an configuration file for the inventory application.

This will create a file named inventory.example.yaml in the current directory.
The file will contain the following sections:
- server.address 	
	* the IP:port of the inventory server
- verbose 		
	* true/false - whether to print verbose output
- location 		
	* the location of the server (freeform text)
- description 		
	* the description of the server (freeform text)

In order to use this configuration automatically, you must move it to one of 
the following locations:

- /etc/inventory/
- ~/.config/inventory/

The file must be named "inventory.yml" to be picked up automatically.

Example:
inventory config
sudo mkdir -p /etc/inventory
sudo mv inventory.example.yaml /etc/inventory/inventory.yaml

Be sure to edit the file to set your actual server address and location.
The server.address is the IP:port of the inventory server.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// check if the config file exists
			_, err := os.Stat(config.GetString("config-file"))
			if err == nil {
				// config file exists, return an error
				cmd.Logger.Error("Config file already exists, please remove it or use a different name")
				return err
			}
			// Apply default values to the config
			config.Set("client.remote", "192.168.5.1:9999")
			config.Set("client.location", "Home Lab")
			config.Set("client.description", "Generic Server")

			config.Set("server.listen", "0.0.0.0")
			config.Set("server.http-port", 8000)
			config.Set("server.rpc-port", 9999)
			svcs := []inventory.Service{
				{
					Name: "syncthing",
					Unit: "syncthing@.service",
					Listeners: []*inventory.Listen{
						{
							Port:          8384,
							ListenAddress: "0.0.0.0",
							Protocol:      "tcp",
						},
						{
							Port:          22000,
							ListenAddress: "0.0.0.0",
							Protocol:      "tcp",
						},
					},
				},
			}
			config.Set("services", svcs)
			config.Set("log-level", 0)
			cfgMap := config.AllSettings()
			// Filter some config settings
			delete(cfgMap, "app")
			delete(cfgMap, "config-file")
			delete(cfgMap, "config-file-used")
			v := viper.New()
			if err := v.MergeConfigMap(cfgMap); err != nil {
				return err
			}

			err = v.WriteConfigAs(config.GetString("config-file"))
			if err != nil {
				cmd.Logger.Error("Error writing config file", "error", err)
				return err
			}
			ui.Success("Config file created", config.GetString("config-file"))
			ui.Info("Config file created:",
				config.GetString("config-file"),
				"Move the file to /etc/inventory/inventory.yml or ~/.config/inventory/inventory.yml",
				"to use it automatically")
			ui.Info("Example:",
				ui.Code("`inventory config`"),
				ui.Code("`sudo mkdir -p /etc/inventory`"),
				ui.Code("`sudo mv inventory.yml /etc/inventory/inventory.yml`"))
			ui.Info("Edit the file to set your actual server address and location.")
			return nil
		},
	}

	return configCmd
}
