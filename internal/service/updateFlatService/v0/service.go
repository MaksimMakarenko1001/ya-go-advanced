package v0

import (
	"context"
	"strconv"

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

func (srv *Service) Do(
	ctx context.Context, metricType, metricName, metricValue string,
) (err error) {
	switch metricType {
	case pkg.MetricTypeCounter:
		if valueInt, err := strconv.ParseInt(metricValue, 10, 64); err != nil {
			return errInvalidMetricValue
		} else {
			return srv.updateCounterService.Do(ctx, metricName, valueInt)
		}

	case pkg.MetricTypeGauge:
		if valueFloat, err := strconv.ParseFloat(metricValue, 64); err != nil {
			return errInvalidMetricValue
		} else {
			return srv.updateGaugeService.Do(ctx, metricName, valueFloat)
		}

	default:
		return errInvalidMetricType
	}
}
