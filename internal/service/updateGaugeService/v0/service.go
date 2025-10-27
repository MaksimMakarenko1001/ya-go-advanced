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
	metricName string, metricValue float64,
) (err error) {
	ok, err := srv.metricRepository.Update(metricName, metricValue)
	if err != nil {
		return pkg.ErrInternalServer.SetInfo(err.Error())
	}
	if !ok {
		return pkg.ErrBadRequest
	}

	return nil
}
