package v0

type MetricRepository interface {
	Update(name string, value float64) (ok bool)
}
