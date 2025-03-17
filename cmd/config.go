/*
Copyright © 2025 Brian Ketelsen <bketelsen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create an example configuration file for the inventory client (reporter)",
	Long: `Create an example configuration file for the inventory client (reporter).

This will create a file named inventory.config.yaml in the current directory.
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
- ~/.inventory/

The file must be named "inventory" with no extension.

Example:
inventory config
sudo mkdir -p /etc/inventory
sudo mv inventory.config.yaml /etc/inventory/inventory

Be sure to edit the file to set your actual server address and location.
The server.address is the IP:port of the inventory server.`,

	Run: func(cmd *cobra.Command, args []string) {

		v := viper.New()
		v.SetConfigName("inventory")
		v.SetConfigType("yaml")

		v.Set("server.address", "192.168.1.10:9999")
		v.Set("location", "Office Rack, Shelf 1")
		v.Set("description", "2U AMD Ryzen 9 5950X")
		v.Set("verbose", false)
		v.WriteConfigAs("inventory.config.yaml")

		cmd.Println(ui.Info("Sample config file created:",
			"./inventory.config.yaml",
			"Move the file to /etc/inventory/inventory or ~/.inventory/inventory",
			"-- be sure to remove the .config.yaml extension --",
			"to use it automatically"))
		cmd.Println(ui.Info("Example:",
			ui.Code("`inventory config`"),
			ui.Code("`sudo mkdir -p /etc/inventory`"),
			ui.Code("`sudo mv inventory.config.yaml /etc/inventory/inventory`")))
		cmd.Println(ui.Info("Edit the file to set your actual server address and location."))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

}
