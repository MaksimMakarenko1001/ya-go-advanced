package updateGaugeV0Service

import (
	"errors"
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
	metricName string, metricValue float64,
) (err error) {
	if ok := srv.metricRepository.Update(metricName, metricValue); !ok {
		return errors.New("error occurred in updateGaugeV0Service.Do")
	}

	return nil
}
