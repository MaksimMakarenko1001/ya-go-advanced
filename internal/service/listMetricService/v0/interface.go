package v0

import (
	"context"
)

type MetricRepository interface {
	List(ctx context.Context) (resp MetricData, err error)
}
