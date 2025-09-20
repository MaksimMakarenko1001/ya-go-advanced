package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/pkg"
)

type (
	UpdateGaugeService   func(metricName string, metricValue float64) (err error)
	UpdateCounterService func(metricName string, metricValue int64) (err error)
)

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
