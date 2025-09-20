package handler

import (
	"net/http"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/pkg"
)

const TypeContentTextPlain = "text/plain"

var allowMetricType = map[string]struct{}{
	pkg.MetricTypeGauge:   {},
	pkg.MetricTypeCounter: {},
}

type Middleware func(next http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func MiddlewareTypeContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if rq.Header.Get("Content-Type") != TypeContentTextPlain {
			http.Error(w, "not supported Content-Type", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, rq)
	})
}

func MiddlewareUrlPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if len(strings.Split(rq.URL.Path, "/")) != 5 {
			http.Error(w, "invalid URL", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, rq)
	})
}

func MiddlewareMetricType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		parts := strings.Split(rq.URL.Path, "/")
		if _, ok := allowMetricType[parts[2]]; !ok {
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, rq)
	})
}

func MiddlewareMetricName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		name := strings.Split(rq.URL.Path, "/")[3]
		if name == "" {
			http.Error(w, "invalid metric name", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, rq)
	})
}
