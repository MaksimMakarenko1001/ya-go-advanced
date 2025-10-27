package v0

import (
	"bytes"
	"context"
	"html/template"
	"slices"
	"strings"

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

func (srv *Service) Do(ctx context.Context, html string) (index string, err error) {
	list, err := srv.metricRepository.List(ctx)
	if err != nil {
		return "", pkg.ErrInternalServer.SetInfo(err.Error())
	}

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
