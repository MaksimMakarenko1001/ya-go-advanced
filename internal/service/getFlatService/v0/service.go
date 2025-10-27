package v0

import (
	"context"
	"strconv"

	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

type Service struct {
	getCounterService *getCounterService.Service
	getGaugeService   *getGaugeService.Service
}

func New(
	getCounterService *getCounterService.Service,
	getGaugeService *getGaugeService.Service,
) *Service {
	return &Service{
		getCounterService: getCounterService,
		getGaugeService:   getGaugeService,
	}
}

func (srv *Service) Do(ctx context.Context, metricType, metricName string) (value string, err error) {
	switch metricType {
	case pkg.MetricTypeCounter:
		if valueInt, err := srv.getCounterService.Do(ctx, metricName); err != nil {
			return "", err
		} else {
			return strconv.FormatInt(*valueInt, 10), nil
		}

	case pkg.MetricTypeGauge:
		if valueFloat, err := srv.getGaugeService.Do(ctx, metricName); err != nil {
			return "", err
		} else {
			return strconv.FormatFloat(*valueFloat, 'f', -1, 64), nil
		}

	default:
		return "", pkg.ErrBadRequest.SetInfo("invalid metric type")
	}
}
