package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
)

type MetricRepository interface {
	GetGauge(ctx context.Context, name string) (item *entities.GaugeItem, ok bool, err error)
}
