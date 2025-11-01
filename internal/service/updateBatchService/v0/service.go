package v0

import (
	"context"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

var (
	errInvalidMetricValue *pkg.Error = pkg.ErrBadRequest.SetInfo("invalid metric value")
	errInvalidMetricType  *pkg.Error = pkg.ErrBadRequest.SetInfo("invalid metric type")
)

type Service struct {
	metricRepository MetricRepository
}

func New(
	metricRepository MetricRepository,
) *Service {
	return &Service{
		metricRepository: metricRepository,
	}
}

func (srv *Service) Do(ctx context.Context, metrics []models.Metrics) (err error) {
	if len(metrics) == 0 {
		return nil
	}

	ts := time.Now()
	seen := make(map[string]struct{}, len(metrics))
	counters := make([]entities.CounterItem, 0, len(metrics))
	gauges := make([]entities.GaugeItem, 0, len(metrics))

	for _, metric := range metrics {
		if _, ok := seen[metric.ID]; ok {
			continue
		}
		seen[metric.ID] = struct{}{}

		switch metric.MType {
		case pkg.MetricTypeCounter:
			if metric.Delta == nil {
				return errInvalidMetricValue
			}
			counters = append(counters, entities.CounterItem{
				MetricType:  metric.MType,
				MetricName:  metric.ID,
				MetricValue: *metric.Delta,
				CreatedAt:   ts,
				UpdatedAt:   ts,
			})

		case pkg.MetricTypeGauge:
			if metric.Value == nil {
				return errInvalidMetricValue
			}
			gauges = append(gauges, entities.GaugeItem{
				MetricType:  metric.MType,
				MetricName:  metric.ID,
				MetricValue: *metric.Value,
				CreatedAt:   ts,
				UpdatedAt:   ts,
			})

		default:
			return errInvalidMetricType
		}
	}

	ok, err := srv.metricRepository.AddUpdateBatch(ctx, counters, gauges)
	if err != nil {
		return pkg.ErrInternalServer.SetInfo(err.Error())
	}
	if !ok {
		return pkg.ErrBadRequest
	}

	return nil
}
