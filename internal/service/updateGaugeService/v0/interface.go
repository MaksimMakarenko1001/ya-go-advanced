package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
)

type MetricRepository interface {
	Update(ctx context.Context, item entities.GaugeItem) (ok bool, err error)
}
