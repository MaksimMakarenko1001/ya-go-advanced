package config

import "flag"

var httpCfg HTTPServerConfig

type diConfig struct {
	HTTP HTTPServerConfig `json:"http"`
}

func (cfg *diConfig) loadConfig() {
	flag.StringVar(&httpCfg.Address, "a", `:8080`, "server net address")

	cfg.HTTP = httpCfg
}

type HTTPServerConfig struct {
	Address string `json:"address"`
}
