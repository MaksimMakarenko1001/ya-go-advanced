package agent

import (
	"flag"
	"time"
)

var (
	httpCfg HTTPClientConfig
	options struct {
		pool   int
		report int
	}
)

type Config struct {
	HTTP           HTTPClientConfig `json:"http"`
	PollInterval   time.Duration    `json:"poll"`
	ReportInterval time.Duration    `json:"report"`
}

func (cfg *Config) LoadConfig() {
	flag.StringVar(&httpCfg.Address, "a", `localhost:8080`, "agent net address")
	flag.IntVar(&options.pool, "p", 2, "pool interval in seconds")
	flag.IntVar(&options.report, "r", 10, "report interval in seconds")

	cfg.HTTP = httpCfg
	cfg.HTTP.Timeout = time.Second * 10

	cfg.PollInterval = time.Second * time.Duration(options.pool)
	cfg.ReportInterval = time.Second * time.Duration(options.report)

}

type HTTPClientConfig struct {
	Address string        `json:"address"`
	Timeout time.Duration `json:"timeout"`
}
