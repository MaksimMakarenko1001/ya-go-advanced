package v0

import (
	"bytes"
	"html/template"
)

type Service struct {
	metricRepository MetricRepository
}

func New(metricRepo MetricRepository) *Service {
	return &Service{
		metricRepository: metricRepo,
	}
}

func (srv *Service) Do(html string) (index string, err error) {
	list := srv.metricRepository.List()

	tmpl := template.Must(template.New("html").Parse(html))

	buffer := bytes.Buffer{}

	err = tmpl.Execute(&buffer, list)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
