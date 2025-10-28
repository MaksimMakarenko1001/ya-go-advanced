package db

import (
	"errors"
	"strconv"
	"strings"
)

type Config struct {
	DSN      string `json:"-"`
	Host     string `env:"host" envDefault:"localhost"`
	Port     uint16 `env:"port" envDefault:"5432"`
	User     string `env:"user" envDefault:"postgres"`
	Password string `env:"password" envDefault:"postgres"`
	Name     string `env:"name" envDefault:"postgres"`
	SSLMode  string `env:"ssl_mode" envDefault:"disable"`
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
	if cfg.Host == "" {
		return errors.New("host undefined")
	}
	if cfg.Port == 0 {
		return errors.New("port undefined")
	}
	if cfg.Name == "" {
		return errors.New("name undefined")
	}
	if cfg.User == "" {
		return errors.New("user undefined")
	}
	return nil
}
