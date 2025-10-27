package v0

import "context"

type MetricRepository interface {
	GetGauge(ctx context.Context, name string) (value float64, ok bool, err error)
}
