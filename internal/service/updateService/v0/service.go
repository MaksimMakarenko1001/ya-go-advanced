package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateCounterService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateGaugeService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

var (
	errInvalidMetricValue *pkg.Error = pkg.ErrBadRequest.SetInfo("invalid metric value")
	errInvalidMetricType  *pkg.Error = pkg.ErrBadRequest.SetInfo("invalid metric type")
)

type Service struct {
	updateCounterService *updateCounterService.Service
	updateGaugeService   *updateGaugeService.Service
}

func New(
	updateCounterService *updateCounterService.Service,
	updateGaugeService *updateGaugeService.Service,
) *Service {
	return &Service{
		updateCounterService: updateCounterService,
		updateGaugeService:   updateGaugeService,
	}
}

func (srv *Service) Do(ctx context.Context, metric models.Metrics) (err error) {
	switch metric.MType {
	case pkg.MetricTypeCounter:
		if metric.Delta == nil {
			return errInvalidMetricValue
		}
		return srv.updateCounterService.Do(ctx, metric.ID, *metric.Delta)

	case pkg.MetricTypeGauge:
		if metric.Value == nil {
			return errInvalidMetricValue
		} else {
			return srv.updateGaugeService.Do(ctx, metric.ID, *metric.Value)
		}

	default:
		return errInvalidMetricType
	}
}
