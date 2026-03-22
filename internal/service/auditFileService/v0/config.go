package v0

type Config struct {
	AuditEnabled bool `env:"audit_enabled" envDefault:"true"`
}
