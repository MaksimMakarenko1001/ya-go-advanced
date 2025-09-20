package updateCounterV0Service

type MetricRepository interface {
	Add(name string, value int64) (ok bool)
}
