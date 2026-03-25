package v0

type Config struct {
	CryptoKey      string `env:"CRYPTO_KEY" json:"cryptoKey"`
	DecryptEnabled bool   `env:"DECRYPT_ENABLED" envDefault:"true" json:"decryptEnabled"`
}
