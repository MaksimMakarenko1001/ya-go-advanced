package v0

import (
	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/pkg"
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
	metricName string,
) (*float64, error) {
	value, ok := srv.metricRepository.Get(metricName)
	if !ok {
		return nil, pkg.ErrNotFound
	}

	metricValue, ok := value.(float64)
	if !ok {
		return nil, pkg.ErrNotFound
	}

	return &metricValue, nil
}
