package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type diConfig struct {
	HTTP HTTPServerConfig `envPrefix:"HTTP_"`
}

func (cfg *diConfig) loadConfig(envPrefix string) {
	cfg.loadFromArg()
	cfg.loadFromEnv(envPrefix)
}

func (cfg *diConfig) loadFromArg() {
	flag.StringVar(&cfg.HTTP.Address, "a", `:8080`, "server net address")

	flag.Parse()
}

func (cfg *diConfig) loadFromEnv(envPrefix string) {
	var config diConfig

	if err := env.Parse(&config, env.Options{Prefix: envPrefix, RequiredIfNoDef: true}); err != nil {
		log.Printf("load from env not ok, %s\n", err.Error())
		return
	}

	if address := config.HTTP.Address; address != "" {
		cfg.HTTP.Address = address
	}
}

type HTTPServerConfig struct {
	Address string `env:"ADDRESS"`
}
