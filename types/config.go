package types

import (
	"github.com/spf13/viper"
)

type Config struct {
	Verbose     bool      `mapstructure:"verbose"`
	Server      Server    `mapstructure:"server"`
	Services    []Service `mapstructure:"services"`
	HTTPPort    int       `mapstructure:"http_port"`
	RPCPort     int       `mapstructure:"rpc_port"`
	Location    string    `mapstructure:"location"`
	Description string    `mapstructure:"description"`
}

type Server struct {
	Address string
}

func ReadConfig() (Config, error) {
	var c Config
	err := viper.Unmarshal(&c)
	return c, err
}
