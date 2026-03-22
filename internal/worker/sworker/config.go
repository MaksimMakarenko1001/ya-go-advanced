package sworker

import "time"

type Config struct {
	JobTimeout  time.Duration `env:"job_timeout" envDefault:"3s"`
	JobInterval time.Duration `env:"job_interval" envDefault:"2s"`
}
