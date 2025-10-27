package v0

type MetricRepository interface {
	Add(name string, value int64) (ok bool, err error)
}
