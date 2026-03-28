package db

import (
	"errors"
	"strconv"
	"strings"
)

type Config struct {
	DSN        string `json:"-"`
	Host       string `env:"HOST" envDefault:"localhost" json:"host"`
	Port       uint16 `env:"PORT" envDefault:"5432" json:"port"`
	User       string `env:"USER" envDefault:"postgres" json:"user"`
	Password   string `env:"PASSWORD" envDefault:"postgres" json:"-"`
	Name       string `env:"NAME" envDefault:"postgres" json:"name"`
	SSLMode    string `env:"SSL_MODE" envDefault:"disable" json:"sslMode"`
	MaxRetries uint16 `env:"MAX_RETRIES" envDefault:"3" json:"maxRetries"`
}

func (cfg Config) ToDSN() (string, error) {
	if err := cfg.validate(); err != nil {
		return "", err
	}

	kv := []string{
		"host=" + cfg.Host,
		"port=" + strconv.Itoa(int(cfg.Port)),
		"user=" + cfg.User,
		"password=" + cfg.Password,
		"dbname=" + cfg.Name,
		"sslmode=" + cfg.SSLMode,
	}

	return strings.Join(kv, " "), nil

}

func (cfg Config) validate() error {
	var errs []error
	if cfg.Host == "" {
		errs = append(errs, errHostUndefined)
	}
	if cfg.Port == 0 {
		errs = append(errs, errPortUndefined)
	}
	if cfg.Name == "" {
		errs = append(errs, errDBNameUndefined)
	}
	if cfg.User == "" {
		errs = append(errs, errUserUndefined)
	}

	return errors.Join(errs...)
}
