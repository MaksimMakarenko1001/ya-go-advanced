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
	metricName string,
) (*int64, error) {
	value, ok := srv.metricRepository.Get(metricName)
	if !ok {
		return nil, pkg.ErrNotFound.SetInfof("`%s` not found", metricName)
	}

	metricValue, ok := value.(int64)
	if !ok {
		return nil, pkg.ErrBadRequest.SetInfof("`%s` type mismatch", metricName)
	}

	return &metricValue, nil
}
