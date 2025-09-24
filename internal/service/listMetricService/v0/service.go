package v0

import (
	"bytes"
	"html/template"
	"slices"
	"strings"
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

	slices.SortFunc(list, func(a, b MetricItem) int {
		return strings.Compare(a.Name, b.Name)
	})

	tmpl := template.Must(template.New("html").Parse(html))

	buffer := bytes.Buffer{}

	err = tmpl.Execute(&buffer, list)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
