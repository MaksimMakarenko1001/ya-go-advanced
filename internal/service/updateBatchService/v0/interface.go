package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
)

type MetricRepository interface {
	AddUpdateBatch(ctx context.Context, ipAddress string, counters []entities.CounterItem, gauges []entities.GaugeItem) (ok bool, err error)
}
