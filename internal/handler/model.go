package handler

import (
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	response *logger.ResponseInfo
}

func (r *responseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.response.Size += size
	return size, err
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.response.Status = statusCode
}
