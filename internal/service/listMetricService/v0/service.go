package v0

import (
	"bytes"
	"html/template"
)

const html = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <table>
			<tbody>
				{{ range . }}
				<tr>
					<td>{{ .Name }}</td>
					<td>{{ .Value }}</td>
				</tr>
				{{ end }}
			</tbody>
		</table>
    </body>
</html>`

type Service struct {
	metricRepository MetricRepository
}

func New(metricRepo MetricRepository) *Service {
	return &Service{
		metricRepository: metricRepo,
	}
}

func (srv *Service) Do() (index string, err error) {
	list := srv.metricRepository.List()

	tmpl := template.Must(template.New("html").Parse(html))

	buffer := bytes.Buffer{}

	err = tmpl.Execute(&buffer, list)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
