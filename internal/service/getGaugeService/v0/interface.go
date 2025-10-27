package v0

type MetricRepository interface {
	GetGauge(name string) (value float64, ok bool, err error)
}
