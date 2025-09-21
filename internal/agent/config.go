package agent

import "time"

type Config struct {
	HTTP           HTTPClientConfig `json:"http"`
	PollInterval   time.Duration    `json:"poll"`
	ReportInterval time.Duration    `json:"report"`
}

func (cfg *Config) LoadConfig() {
	cfg.HTTP = HTTPClientConfig{
		Host:    `localhost:8080`,
		Timeout: time.Second * 10,
	}
	cfg.PollInterval = time.Second * 2
	cfg.ReportInterval = time.Second * 10
}

type HTTPClientConfig struct {
	Host    string        `json:"host"`
	Timeout time.Duration `json:"timeout"`
}
