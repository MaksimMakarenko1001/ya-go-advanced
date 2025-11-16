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
	Address        string        `env:"ADDRESS"`
	Timeout        time.Duration `env:"TIMEOUT" envDefault:"10s"`
	BatchSize      int           `env:"BATCH_SIZE" envDefault:"3"`
	MaxRetries     uint16        `env:"MAX_RETRIES" envDefault:"3"`
	Key            string        `env:"KEY"`
	PollInterval   time.Duration `env:"POOL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

func (cfg *Config) LoadConfig(envPrefix string) {
	cfg.loadFromArg()
	cfg.loadFromEnv(envPrefix)
	cfg.loadFromEnvPassTests() // meets tests required
}

func (cfg *Config) loadFromArg() {
	flag.StringVar(&cfg.Address, "a", `localhost:8080`, "agent net address")
	flag.DurationVar(&cfg.PollInterval, "p", 2, "pool interval in seconds")
	flag.DurationVar(&cfg.ReportInterval, "r", 10, "report interval in seconds")
	flag.StringVar(&cfg.Key, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", 5, "num threads work concurrently")

	flag.Parse()

	cfg.PollInterval = time.Second * cfg.PollInterval
	cfg.ReportInterval = time.Second * cfg.ReportInterval

}

func (cfg *Config) loadFromEnv(envPrefix string) {
	var config Config

	if err := env.Parse(&config, env.Options{Prefix: envPrefix}); err != nil {
		log.Printf("env not ok, %s\n", err.Error())
		return
	}

	cfg.BatchSize = config.BatchSize
	cfg.MaxRetries = config.MaxRetries
	cfg.Key = config.Key

	if address := config.Address; address != "" {
		cfg.Address = address
	}
	if pool := config.PollInterval; pool.String() != "0s" {
		cfg.PollInterval = pool
	}
	if report := config.ReportInterval; report.String() != "0s" {
		cfg.ReportInterval = report
	}
	if rateLimit := config.RateLimit; rateLimit > 0 {
		cfg.RateLimit = rateLimit
	}
}

func (cfg *Config) loadFromEnvPassTests() {
	if address := os.Getenv("ADDRESS"); address != "" {
		cfg.Address = address
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

	if key := os.Getenv("KEY"); key != "" {
		cfg.Key = key
	}

	if rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT")); err != nil {
		log.Printf("RATE_LIMIT env not ok, %s\n", err.Error())
	} else {
		cfg.RateLimit = rateLimit
	}
}
