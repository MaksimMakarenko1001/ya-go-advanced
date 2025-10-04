package agent

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	HTTP           HTTPClientConfig `envPrefix:"HTTP_"`
	PollInterval   time.Duration    `env:"POOL_INTERVAL"`
	ReportInterval time.Duration    `env:"REPORT_INTERVAL"`
}

func (cfg *Config) LoadConfig(envPrefix string) {
	cfg.loadFromArg()
	cfg.loadFromEnv(envPrefix)
}

func (cfg *Config) loadFromArg() {
	var options struct {
		pool   int
		report int
	}

	flag.StringVar(&cfg.HTTP.Address, "a", `localhost:8080`, "agent net address")
	flag.IntVar(&options.pool, "p", 2, "pool interval in seconds")
	flag.IntVar(&options.report, "r", 10, "report interval in seconds")

	flag.Parse()

	cfg.PollInterval = time.Second * time.Duration(options.pool)
	cfg.ReportInterval = time.Second * time.Duration(options.report)
}

func (cfg *Config) loadFromEnv(envPrefix string) {
	var config Config

	if err := env.Parse(&config, env.Options{Prefix: envPrefix}); err != nil {
		log.Printf("load from env not ok, %s\n", err.Error())
		return
	}

	if address := config.HTTP.Address; address != "" {
		cfg.HTTP.Address = address
	}
	if pool := config.PollInterval; pool.String() != "0s" {
		cfg.PollInterval = pool
	}
	if report := config.ReportInterval; report.String() != "0s" {
		cfg.ReportInterval = report
	}
}

type HTTPClientConfig struct {
	Address string        `env:"ADDRESS"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"10s"`
}
