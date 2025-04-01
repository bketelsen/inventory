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
	"log/slog"
	"os"
	"strings"

	"github.com/bketelsen/toolbox/cobra"
	goversion "github.com/bketelsen/toolbox/go-version"
	"github.com/bketelsen/toolbox/slug"
	"github.com/spf13/viper"
)

var cfgFile string
var foundConfig bool
var appname = "inventory"
var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

var bversion = buildVersion(version, commit, date, builtBy, treeState)

// ldflags
// Default: '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser'.

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "inventory",
	Version: bversion.String(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			cmd.GlobalConfig().SetConfigFile(cfgFile) // Use config file from the flag if set
		} else {
			cmd.GlobalConfig().SetConfigName("inventory.yaml") // name of config file

		}
		if err := cmd.GlobalConfig().ReadInConfig(); err == nil {
			defer slog.Debug("Using config file:", slog.String("file", cmd.GlobalConfig().ConfigFileUsed()))
			foundConfig = true
		} else {
			defer slog.Debug("Error reading config file", slug.Err(err))
			foundConfig = false
		}
		// set the slog default logger to the cobra logger
		slog.SetDefault(cmd.Logger)
		// set log level based on the --verbose flag
		if cmd.GlobalConfig().GetBool("verbose") {
			cmd.SetLogLevel(slog.LevelDebug)
			cmd.Logger.Debug("Debug logging enabled")
		}
	},
	InitConfig: func() *viper.Viper {
		config := viper.New()
		config.SetEnvPrefix(appname)
		config.AutomaticEnv()
		config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
		config.SetConfigType("yaml")
		config.SetConfigName("inventory.yaml")          // name of config file
		config.AddConfigPath("/etc/inventory/")         // path to look for the config file in
		config.AddConfigPath("$HOME/.config/inventory") // call multiple times to add many search paths
		config.AddConfigPath("$HOME/.inventory")        // deprecate this in the future, keep home clean
		config.AddConfigPath(".")                       // optionally look for config in the working directory

		return config
	},
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/inventory/inventory.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")

}

// https://www.asciiart.eu/text-to-ascii-art to make your own
// just make sure the font doesn't have backticks in the letters or
// it will break the string quoting
var asciiName = `
██╗███╗   ██╗██╗   ██╗███████╗███╗   ██╗████████╗ ██████╗ ██████╗ ██╗   ██╗
██║████╗  ██║██║   ██║██╔════╝████╗  ██║╚══██╔══╝██╔═══██╗██╔══██╗╚██╗ ██╔╝
██║██╔██╗ ██║██║   ██║█████╗  ██╔██╗ ██║   ██║   ██║   ██║██████╔╝ ╚████╔╝ 
██║██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██║╚██╗██║   ██║   ██║   ██║██╔══██╗  ╚██╔╝  
██║██║ ╚████║ ╚████╔╝ ███████╗██║ ╚████║   ██║   ╚██████╔╝██║  ██║   ██║   
╚═╝╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   
`

// buildVersion builds the version info for the application
func buildVersion(version, commit, date, builtBy, treeState string) goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails(appname, "Collect and report deployment information.", "https://bketelsen.github.io/inventory"),
		goversion.WithASCIIName(asciiName),
		func(i *goversion.Info) {
			if commit != "" {
				i.GitCommit = commit
			}
			if treeState != "" {
				i.GitTreeState = treeState
			}
			if date != "" {
				i.BuildDate = date
			}
			if version != "" {
				i.GitVersion = version
			}
			if builtBy != "" {
				i.BuiltBy = builtBy
			}

		},
	)
}
