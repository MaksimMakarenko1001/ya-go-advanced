package v0

import "context"

type MetricRepository interface {
	GetCounter(ctx context.Context, name string) (value int64, ok bool, err error)
}
