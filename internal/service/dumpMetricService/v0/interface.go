package v0

type MetricRepository interface {
	Load(b []byte) error
	Save() ([]byte, error)
}
