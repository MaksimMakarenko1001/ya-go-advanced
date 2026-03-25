package sworker

import "time"

type Config struct {
	JobTimeout  time.Duration `env:"JOB_TIMEOUT" envDefault:"3s" json:"jobTimeout"`
	JobInterval time.Duration `env:"JOB_INTERVAL" envDefault:"2s" json:"jobInterval"`
}
