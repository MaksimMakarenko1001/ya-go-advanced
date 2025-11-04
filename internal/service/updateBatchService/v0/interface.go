package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
)

type MetricRepository interface {
	AddUpdateBatch(ctx context.Context, counters []entities.CounterItem, gauges []entities.GaugeItem) (ok bool, err error)
}
