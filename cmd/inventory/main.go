package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/bketelsen/toolbox/cobra"
	"github.com/bketelsen/toolbox/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"go.uber.org/automaxprocs/maxprocs"

	goversion "github.com/bketelsen/toolbox/go-version"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"
)

var (
	// The global config object
	logLevel slog.Level
	appname  = "inventory"

	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

// LogLevelIDs Maps 3rd party enumeration values to their textual representations
var LogLevelIDs = map[slog.Level][]string{
	slog.LevelDebug: {"debug"},
	slog.LevelInfo:  {"info"},
	slog.LevelWarn:  {"warn"},
	slog.LevelError: {"error"},
}

func main() {
	cmd, config := NewRootCommand()
	cmd.AddCommand(NewServerCommand(config))
	cmd.AddCommand(NewSendCommand(config))
	cmd.AddCommand(NewConfigCommand(config))
	cmd.AddCommand(NewSearchCommand(config))
	cmd.AddCommand(NewManCommand(config))
	cmd.AddCommand(NewGendocsCommand(config))
	cmd.AddCommand(NewChangelogCommand(config))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ldflags
// Default: '-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}} -X main.builtBy=goreleaser'.
var bversion = buildVersion(version, commit, date, builtBy, treeState)

// NewRootCommand creates a new root command for the application
func NewRootCommand() (*cobra.Command, *viper.Viper) {
	// this sets the config file location & file name to the current working directory
	// and the name of the package
	// this is the default location for the config file
	cwd, _ := os.Getwd()
	cfgFile := path.Join(cwd, fmt.Sprintf("%s.yml", appname))

	config := setupConfig()

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "inventory",
		Short: "Inventory is a tool to collect and report deployment information",
		Long: ui.Long("Inventory is a tool to collect and report deployment information to a central server. It collects information about the host, docker/incus containers, and manually specified services running on the host. The reporting command is designed to be run as a cron job or systemd timer.",
			ui.Example{
				Description: "Send inventory to the server",
				Command:     "inventory send",
			},
			ui.Example{
				Description: "Send inventory to the server with debug logging",
				Command:     "inventory send --log-level debug",
			},
			ui.Example{
				Description: "Send inventory to the server with a custom config file",
				Command:     "inventory send --config-file /path/to/config.yaml",
			},
			ui.Example{
				Description: "Start the server",
				Command:     "inventory server",
			},
			ui.Example{
				Description: "Start the server with debug logging",
				Command:     "inventory server --log-level debug",
			},
			ui.Example{
				Description: "Start the server with a custom config file",
				Command:     "inventory server --config-file /path/to/config.yaml",
			},
		),
		Version: bversion.String(),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default slog logger to the cobra command
			slog.SetDefault(cmd.Logger)

			// bind the slog log level to our enumflag
			ll := config.GetInt("log-level")
			cmd.SetLogLevel(slog.Level(ll))

			// only prints if the log level is set to debug
			cmd.Logger.Debug("Debug logging enabled")

			// if the pflag has a value other than the default, then
			// reload the config file
			if cmd.Flags().Lookup("config-file").Changed {
				slog.Debug("Using config file from flag", "file", cfgFile)
				config.SetConfigFile(cfgFile)
				config.Set("config-file", cfgFile)
				return config.ReadInConfig()
			}
			// otherwise use the default config
			// created or loaded by the setupConfig function
			slog.Info("Config Used", "file", config.ConfigFileUsed())

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config-file",
		"c",
		cfgFile,
		``)
	_ = config.BindPFlag("config-file", rootCmd.PersistentFlags().Lookup("config-file"))
	rootCmd.PersistentFlags().Var(
		enumflag.New(&logLevel, "log", LogLevelIDs, enumflag.EnumCaseInsensitive),
		"log-level",
		"logging level [debug|info|warn|error]")
	return rootCmd, config
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

func init() {
	// enable colored output on github actions et al
	if os.Getenv("NOCOLOR") != "" {
		lipgloss.DefaultRenderer().SetColorProfile(termenv.Ascii)
	}
	// automatically set GOMAXPROCS to match available CPUs.
	// GOMAXPROCS will be used as the default value for the --parallelism flag.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("failed to set GOMAXPROCS")
	}
}
