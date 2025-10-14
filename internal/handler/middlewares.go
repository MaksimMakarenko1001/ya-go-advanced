package handler

import (
	"compress/gzip"
	"net/http"
	"slices"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

const (
	TypeContentTextPlain       = "text/plain"
	TypeContentApplicationJSON = "application/json"
)

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

func MiddlewareTypeContentTextPlain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if rq.Header.Get("Content-Type") != TypeContentTextPlain {
			http.Error(w, "not supported Content-Type", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, rq)
	})
}

func MiddlewareTypeContentApplicationJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if rq.Header.Get("Content-Type") != TypeContentApplicationJSON {
			http.Error(w, "not supported Content-Type", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, rq)
	})
}

func MiddlewareURLPath(next http.Handler) http.Handler {
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

func MiddlewareCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := rw

		supportsGzip := slices.Contains(r.Header.Values("Accept-Encoding"), "gzip")
		if supportsGzip {
			cw := &compressWriter{
				w:  rw,
				zw: gzip.NewWriter(rw),
			}
			w = cw
			defer cw.Close()
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			cr := &compressReader{
				r:  r.Body,
				zr: zr,
			}

			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(w, r)
	})
}
