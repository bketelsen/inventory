/*
Copyright © 2025 Brian Ketelsen

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
	"log"
	"log/slog"

	"github.com/bketelsen/inventory/client"
	"github.com/bketelsen/inventory/types"
	"github.com/bketelsen/toolbox/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
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
		if !foundConfig {
			cmd.Logger.Error("No config file found, cannot send inventory")
			cmd.Logger.Info("Run 'inventory config' to create an example config file")
			return nil
		}
		slog.SetDefault(cmd.Logger)
		cmd.Logger.Info(cmd.GlobalConfig().GetString("description"))
		// Read config
		config, err := types.ViperToStruct(cmd.GlobalConfig())
		if err != nil {
			log.Println("Error reading config:", err)
			return err
		}

		// Get server address from config or use default

		cl := client.NewClient(config)
		err = cl.Send()
		if err != nil {
			log.Printf("Error sending report: %v", err)
			return err
		}
		cmd.Logger.Info("Inventory sent successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
