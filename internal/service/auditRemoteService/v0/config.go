package v0

type Config struct {
	AuditEnabled bool `env:"AUDIT_ENABLED" envDefault:"true" json:"auditEnabled"`
}
