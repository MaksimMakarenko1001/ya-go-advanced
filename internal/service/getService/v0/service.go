package v0

import (
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

var (
	errInvalidMetricType *pkg.Error = pkg.ErrBadRequest.SetInfo("invalid metric type")
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

func (srv *Service) Do(metricType, metricName string) (metric *models.Metrics, err error) {
	resp := models.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case pkg.MetricTypeCounter:
		if valueInt, err := srv.getCounterService.Do(metricName); err != nil {
			return nil, err
		} else {
			resp.Delta = valueInt
		}

	case pkg.MetricTypeGauge:
		if valueFloat, err := srv.getGaugeService.Do(metricName); err != nil {
			return nil, err
		} else {
			resp.Value = valueFloat
		}

	default:
		return nil, errInvalidMetricType
	}

	return &resp, nil
}
