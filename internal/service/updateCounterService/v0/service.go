package v0

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
	metricName string, metricValue int64,
) (err error) {
	if ok := srv.metricRepository.Add(metricName, metricValue); !ok {
		return errors.New("error occurred in updateCounterV0Service.Do")
	}

	return nil
}
