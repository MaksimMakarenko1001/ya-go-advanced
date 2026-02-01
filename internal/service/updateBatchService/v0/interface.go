package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
)

type MetricRepository interface {
	AddUpdateBatch(ctx context.Context, counters []entities.CounterItem, gauges []entities.GaugeItem,
		outboxes []entities.Outbox, outboxSegment string,
	) (ok bool, err error)
}
