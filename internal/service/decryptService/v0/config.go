package v0

type Config struct {
	CryptoKey      string `env:"CRYPTO_KEY"`
	DecryptEnabled bool   `env:"decrypt_enabled" envDefault:"true"`
}
