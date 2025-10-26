package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

const html = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <table>
			<tbody>{{ range . }}
				<tr>
					<td>{{ .Name }}</td>
					<td>{{ .Value }}</td>
				</tr>{{ end }}
			</tbody>
		</table>
    </body>
</html>`

type (
	UpdateFlatService func(metricType, metricName, metricValue string) (err error)
	UpdateService     func(metric models.Metrics) (err error)

	GetGaugeService   func(metricName string) (metricValue *float64, err error)
	GetCounterService func(metricName string) (metricValue *int64, err error)
	GetFlatService    func(metricType, metricName string) (metricValue string, err error)

	ListMetricService func(template string) (index string, err error)
)

func DoListMetricResponse(srv ListMetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		index, err := srv(html)
		if err != nil {
			WriteError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, index)
	}
}

func DoUpdateFlatResponse(srv UpdateFlatService, metricType, metricName, metricValue string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := srv(metricType, metricName, metricValue); err != nil {
			WriteError(w, err)
			return
		}

		WriteOK(w)
	}
}

func DoUpdateJSONResponse(srv UpdateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err := srv(metric); err != nil {
			WriteError(w, err)
			return
		}

		resp, _ := json.Marshal(metric)
		WriteJSONResult(w, resp)
	}
}

func DoGetFlatResponse(srv GetFlatService, metricType, metricName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := srv(metricType, metricName)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, value)
	}
}

func DoGetGaugeJSONResponse(srv GetGaugeService, rq models.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := srv(rq.ID)
		if err != nil {
			WriteError(w, err)
			return
		}

		resp, err := json.Marshal(models.Metrics{
			ID:    rq.ID,
			MType: rq.MType,
			Value: value,
		})
		if err != nil {
			WriteError(w, fmt.Errorf("convert to gauge response not ok, %w", err))
		}

		WriteJSONResult(w, resp)
	}
}

func DoGetCounterJSONResponse(srv GetCounterService, rq models.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := srv(rq.ID)
		if err != nil {
			WriteError(w, err)
			return
		}

		resp, err := json.Marshal(models.Metrics{
			ID:    rq.ID,
			MType: rq.MType,
			Delta: value,
		})
		if err != nil {
			WriteError(w, fmt.Errorf("convert to counter response not ok, %w", err))
		}

		WriteJSONResult(w, resp)
	}
}

func WriteJSONResult(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		WriteError(w, fmt.Errorf("write json not ok, %w", err))
	}
}

func WriteResult(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, res)
}

func WriteOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		err = pkg.ErrInternalServer
	}

	errE, ok := err.(*pkg.Error)
	if !ok {
		errE = pkg.ErrInternalServer
	}
	http.Error(w, errE.Error(), errE.HTTPStatus())
}
