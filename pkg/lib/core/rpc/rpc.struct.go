package rpc

import (
	"github.com/LoveCatdd/util/pkg/lib/core/viper"
)

type RPCConfig struct {
	Rpc struct {
		Server struct {
			Addr string `mapstructure:"addr"`
			Name string `mapstructure:"name"`
		} `mapstructure:"server"`
		Unauth struct {
			Path []string `mapstructure:"path"`
		} `mapstructure:"unauth"`
		Timeout int32      `mapstructure:"timeout"`
		Handler []Handlers `mapstructure:"handler"`
	} `mapstructure:"rpc"`
}

type Handlers struct {
	Type   string `mapstructure:"type"`
	Method string `mapstructure:"method"`
}

func (*RPCConfig) FileType() string {
	return viper.VIPER_YAML
}

var RpcConf = new(RPCConfig)
