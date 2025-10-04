package v0

import (
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
	metricName string, metricValue int64,
) (err error) {
	if ok := srv.metricRepository.Add(metricName, metricValue); !ok {
		return pkg.ErrBadRequest
	}

	return nil
}
