package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
)

type MetricRepository interface {
	GetCounter(ctx context.Context, name string) (item *entities.CounterItem, ok bool, err error)
}
