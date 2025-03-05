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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "inventory",
	Example: `//client
inventory send
inventory send --verbose // more verbose output
inventory send --verbose --config /path/to/config.yaml // specify a config file

//server
inventory server
inventory server --verbose // more verbose output
inventory server --verbose --config /path/to/config.yaml // specify a config file`,
	Short: "Inventory is a tool to collect and report deployment information",
	Long: `Inventory is a tool to collect and report deployment information
to a central server. It collects information about the host,
docker/incus containers, and manually specified services running on the host.
The reporting command is designed to be run as a cron job or systemd timer.

It is not intended to be run interactively.



Inventory listens for http requests on port 8000 by default.
Inventory listens for rpc requests on port 9000 by default.
  `,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/inventory/inventory)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		viper.SetConfigName("inventory")        // name of config file (without extension)
		viper.SetConfigType("yaml")             // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/inventory/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.inventory") // call multiple times to add many search paths
		viper.AddConfigPath(".")                // optionally look for config in the working directory
		err := viper.ReadInConfig()             // Find and read the config file
		if err != nil {                         // Handle errors reading the config file
			if strings.Contains(err.Error(), "control characters") {
				log.Println("No config file found, using defaults")
			} else {
				log.Println(fmt.Errorf("fatal error config file: %w", err))

			}
		} else {
			if viper.GetBool("verbose") {
				log.Println("Using config file:", viper.ConfigFileUsed())
			}
		}
	}
	viper.AutomaticEnv() // read in environment variables that match

}
