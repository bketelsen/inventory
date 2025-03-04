package types

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Verbose  bool      `mapstructure:"verbose"`
	Server   Server    `mapstructure:"server"`
	Services []Service `mapstructure:"services"`
	HTTPPort int       `mapstructure:"http_port"`
	RPCPort  int       `mapstructure:"rpc_port"`
}

type Server struct {
	Address string
}

func ReadConfig() (Config, error) {
	var c Config
	log.Println(viper.ConfigFileUsed())
	err := viper.Unmarshal(&c)
	return c, err
}
