package config

import (
	"flag"
	"log"
	"os"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	"github.com/caarlos0/env/v6"
)

type diConfig struct {
	HTTP   HTTPServerConfig `envPrefix:"HTTP_"`
	Logger logger.Config    `envPrefix:"LOGGER_"`
}

func (cfg *diConfig) loadConfig(envPrefix string) {
	cfg.loadFromArg()
	cfg.loadFromEnv(envPrefix)
	cfg.loadFromEnvToPassTests() // meets tests required
}

func (cfg *diConfig) loadFromArg() {
	flag.StringVar(&cfg.HTTP.Address, "a", `:8080`, "server net address")

	flag.Parse()
}

func (cfg *diConfig) loadFromEnv(envPrefix string) {
	var config diConfig

	if err := env.Parse(&config, env.Options{Prefix: envPrefix, RequiredIfNoDef: true}); err != nil {
		log.Printf("env not ok, %s\n", err.Error())
		return
	}

	if address := config.HTTP.Address; address != "" {
		cfg.HTTP.Address = address
	}
}

func (cfg *diConfig) loadFromEnvToPassTests() {
	if address := os.Getenv("ADDRESS"); address != "" {
		cfg.HTTP.Address = address
	}
}

type HTTPServerConfig struct {
	Address string `env:"ADDRESS"`
}
