package handler

import (
	"context"
	"encoding/json"
	"errors"
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
	UpdateFlatService  func(ctx context.Context, metricType, metricName, metricValue string) (err error)
	UpdateBatchService func(ctx context.Context, metrics []models.Metrics) (err error)
	UpdateService      func(ctx context.Context, metric models.Metrics) (err error)

	GetGaugeService   func(ctx context.Context, metricName string) (metricValue *float64, err error)
	GetCounterService func(ctx context.Context, metricName string) (metricValue *int64, err error)
	GetFlatService    func(ctx context.Context, metricType, metricName string) (metricValue string, err error)
	GetService        func(ctx context.Context, metricType, metricName string) (metric *models.Metrics, err error)

	ListMetricService func(ctx context.Context, template string) (index string, err error)
)

func DoListMetricResponse(srv ListMetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		index, err := srv(r.Context(), html)
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
		if err := srv(r.Context(), metricType, metricName, metricValue); err != nil {
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
			return
		}

		if err := srv(r.Context(), metric); err != nil {
			WriteError(w, err)
			return
		}

		resp, _ := json.Marshal(metric)
		WriteJSONResult(w, resp)
	}
}

func DoUpdateBatchJSONResponse(srv UpdateBatchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := srv(r.Context(), metrics); err != nil {
			WriteError(w, err)
			return
		}

		resp, _ := json.Marshal(metrics)
		WriteJSONResult(w, resp)
	}
}

func DoGetFlatResponse(srv GetFlatService, metricType, metricName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := srv(r.Context(), metricType, metricName)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, value)
	}
}

func DoGetJSONResponse(srv GetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		metric, err := srv(r.Context(), request.MType, request.ID)
		if err != nil {
			WriteError(w, err)
			return
		}

		resp, err := json.Marshal(*metric)
		if err != nil {
			WriteError(w, fmt.Errorf("convert to get response not ok, %w", err))
			return
		}

		WriteJSONResult(w, resp)
	}
}

func WriteJSONResult(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		WriteError(w, fmt.Errorf("write json not ok, %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func WriteResult(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, res)
	w.WriteHeader(http.StatusOK)
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

	var errE *pkg.Error

	if !errors.As(err, &errE) {
		errE = pkg.ErrInternalServer
	}
	http.Error(w, errE.Error(), errE.HTTPStatus())
}
