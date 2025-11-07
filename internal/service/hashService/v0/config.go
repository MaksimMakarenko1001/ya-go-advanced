package v0

type Config struct {
	Key string `env:"KEY" envDefault:"secret"`
}
