package v0

type MetricRepository interface {
	GetCounter(name string) (value int64, ok bool, err error)
}
