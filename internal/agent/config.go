package agent

import (
	"flag"
	"log"
	"os"
	"strconv"
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
	cfg.loadFromEnvPassTests() // meets tests required
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
		log.Printf("env not ok, %s\n", err.Error())
		return
	}

	cfg.HTTP.BatchSize = config.HTTP.BatchSize
	cfg.HTTP.MaxRetries = config.HTTP.MaxRetries

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

func (cfg *Config) loadFromEnvPassTests() {
	if address := os.Getenv("ADDRESS"); address != "" {
		cfg.HTTP.Address = address
	}

	if pool, err := strconv.Atoi(os.Getenv("POOL_INTERVAL")); err != nil {
		log.Printf("POOL_INTERVAL env not ok, %s\n", err.Error())
	} else {
		cfg.PollInterval = time.Second * time.Duration(pool)
	}

	if report, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL")); err != nil {
		log.Printf("REPORT_INTERVAL env not ok, %s\n", err.Error())
	} else {
		cfg.ReportInterval = time.Second * time.Duration(report)
	}
}

type HTTPClientConfig struct {
	Address    string        `env:"ADDRESS"`
	Timeout    time.Duration `env:"TIMEOUT" envDefault:"10s"`
	BatchSize  int           `env:"BATCH_SIZE" envDefault:"3"`
	MaxRetries uint16        `env:"MAX_RETRIES" envDefault:"3"`
}
