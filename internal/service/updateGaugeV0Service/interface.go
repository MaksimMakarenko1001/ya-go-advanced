package updateGaugeV0Service

type MetricRepository interface {
	Update(name string, value float64) (ok bool)
}
