package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	"github.com/caarlos0/env/v6"
)

type diConfig struct {
	HTTP            HTTPServerConfig `envPrefix:"HTTP_"`
	Logger          logger.Config    `envPrefix:"LOGGER_"`
	StoreInterval   time.Duration    `env:"STORE_INTERVAL"`
	FileStoragePath string           `env:"FILE_STORAGE_PATH"`
	Restore         bool             `env:"RESTORE"`
}

func (cfg *diConfig) loadConfig(envPrefix string) {
	cfg.loadFromArg()
	cfg.loadFromEnv(envPrefix)
	cfg.loadFromEnvToPassTests() // meets tests required
}

func (cfg *diConfig) loadFromArg() {
	var options struct {
		store int
	}

	flag.StringVar(&cfg.HTTP.Address, "a", `:8080`, "server net address")
	flag.IntVar(&options.store, "i", 0, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", "dump.txt", "dump file path")
	flag.BoolVar(&cfg.Restore, "r", false, "restore dump file on start")

	flag.Parse()

	cfg.StoreInterval = time.Second * time.Duration(options.store)
}

func (cfg *diConfig) loadFromEnv(envPrefix string) {
	var config diConfig

	if err := env.Parse(&config, env.Options{Prefix: envPrefix /*, RequiredIfNoDef: true*/}); err != nil {
		log.Printf("env not ok, %s\n", err.Error())
		return
	}

	cfg.Logger = config.Logger

	if address := config.HTTP.Address; address != "" {
		cfg.HTTP.Address = address
	}
	if store := config.StoreInterval; store.String() != "0s" {
		cfg.StoreInterval = store
	}
	if fname := config.FileStoragePath; fname != "" {
		cfg.FileStoragePath = fname
	}
	if restore := config.Restore; restore {
		cfg.Restore = restore
	}
}

func (cfg *diConfig) loadFromEnvToPassTests() {
	if address := os.Getenv("ADDRESS"); address != "" {
		cfg.HTTP.Address = address
	}
	if store, err := strconv.Atoi(os.Getenv("STORE_INTERVAL")); err != nil {
		log.Printf("STORE_INTERVAL env not ok, %s\n", err.Error())
	} else {
		cfg.StoreInterval = time.Second * time.Duration(store)
	}
	if fname := os.Getenv("FILE_STORAGE_PATH"); fname != "" {
		cfg.FileStoragePath = fname
	}
	if restore := os.Getenv("RESTORE"); restore != "" {
		cfg.Restore = restore == "true"
	}
}

type HTTPServerConfig struct {
	Address string `env:"ADDRESS"`
}
