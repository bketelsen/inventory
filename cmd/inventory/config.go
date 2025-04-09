package main

import (
	"strings"

	"github.com/spf13/viper"
)

// setupConfig initializes the viper configuration with defaults and environment variables
// It sets the config file name and paths to search for the config file.
// It also sets up the environment variable prefix and key replacer for environment variables.
func setupConfig() *viper.Viper {
	config := viper.New()
	config.SetEnvPrefix(appname)
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
	config.SetDefault("debug", false)
	config.SetConfigType("yaml")
	config.SetConfigName("inventory.yml")           // name of config file
	config.AddConfigPath("/etc/inventory/")         // path to look for the config file in
	config.AddConfigPath("$HOME/.config/inventory") // call multiple times to add many search paths
	config.AddConfigPath(".")                       // optionally look for config in the working directory

	_ = config.ReadInConfig()

	config.Set("app.name", appname)
	config.Set("app.version", version)
	config.Set("app.commit", commit)

	return config
}
