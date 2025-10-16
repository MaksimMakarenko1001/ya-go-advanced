package v0

type Service struct {
	metricRepository MetricRepository
}

func New(metricRepo MetricRepository) *Service {
	return &Service{
		metricRepository: metricRepo,
	}
}

func (srv *Service) ReadFile(name string) error {
	return nil
}

func (srv *Service) WriteFile(name string) error {
	return nil
}
