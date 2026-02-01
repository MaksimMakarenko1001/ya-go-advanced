package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
)

type MetricRepository interface {
	Add(ctx context.Context, item entities.CounterItem) (ok bool, err error)
}
