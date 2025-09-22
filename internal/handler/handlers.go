package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/pkg"
)

type Table struct {
	Content string
}

type (
	UpdateGaugeService   func(metricName string, metricValue float64) (err error)
	UpdateCounterService func(metricName string, metricValue int64) (err error)

	GetGaugeService   func(metricName string) (metricValue *float64, err error)
	GetCounterService func(metricName string) (metricValue *int64, err error)

	ListMetricService func() (template string, err error)
)

func DoListMetricResponse(srv ListMetricService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		index, err := srv()
		if err != nil {
			WriteError(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, index)
	}
}

func DoUpdateGaugeResponse(srv UpdateGaugeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		value, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			http.Error(w, "invalid metric value", http.StatusBadRequest)
			return
		}

		name := parts[3]

		if err := srv(name, value); err != nil {
			WriteError(w, err)
			return
		}

		WriteOK(w)
	}
}

func DoGetGaugeResponse(srv GetGaugeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		name := parts[3]

		value, err := srv(name)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, strconv.FormatFloat(*value, 'f', -1, 64))
	}
}

func DoUpdateCounterResponse(srv UpdateCounterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		value, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			http.Error(w, "invalid metric value", http.StatusBadRequest)
			return
		}

		name := parts[3]

		if err := srv(name, value); err != nil {
			WriteError(w, err)
			return
		}

		WriteOK(w)
	}
}

func DoGetCounterResponse(srv GetCounterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		name := parts[3]

		value, err := srv(name)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, strconv.FormatInt(*value, 10))
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
	if err == nil {
		err = pkg.ErrInternalServer
	}

	errE, ok := err.(*pkg.Error)
	if !ok {
		log.Println(err.Error())
		errE = pkg.ErrInternalServer
	}

	w.WriteHeader(errE.HTTPStatus())

}
