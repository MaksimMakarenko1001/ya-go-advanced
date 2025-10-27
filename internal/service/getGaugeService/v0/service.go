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
) (*float64, error) {
	value, ok, err := srv.metricRepository.GetGauge(metricName)
	if err != nil {
		return nil, pkg.ErrInternalServer.SetInfo(err.Error())
	}
	if !ok {
		return nil, pkg.ErrNotFound.SetInfof("`%s` not found", metricName)
	}

	return &value, nil
}
