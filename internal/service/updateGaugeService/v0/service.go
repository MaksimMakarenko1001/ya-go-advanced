package v0

import (
	"context"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

type Service struct {
	metricRepository MetricRepository
}

func New(metricRepo MetricRepository) *Service {
	return &Service{
		metricRepository: metricRepo,
	}
}

func (srv *Service) Do(
	ctx context.Context, metricName string, metricValue float64,
) (err error) {
	ts := time.Now()

	ok, err := srv.metricRepository.Update(ctx, entities.GaugeItem{
		MetricType:  pkg.MetricTypeGauge,
		MetricName:  metricName,
		MetricValue: metricValue,
		CreatedAt:   ts,
		UpdatedAt:   ts,
	})
	if err != nil {
		return pkg.ErrInternalServer.SetInfo(err.Error())
	}
	if !ok {
		return pkg.ErrBadRequest
	}

	return nil
}
