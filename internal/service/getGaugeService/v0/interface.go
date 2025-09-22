package v0

type MetricRepository interface {
	Get(name string) (value any, ok bool)
}
