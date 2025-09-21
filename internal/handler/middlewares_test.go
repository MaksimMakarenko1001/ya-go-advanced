package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func OKHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func TestMiddlewareTypeContent(t *testing.T) {
	type want struct {
		code    int
		message string
	}
	tests := []struct {
		name        string
		contentType string
		want        want
	}{
		{
			name:        "positive test",
			contentType: "text/plain",
			want: want{
				code:    200,
				message: "",
			},
		},
		{
			name:        "negative test",
			contentType: "application/json",
			want: want{
				code:    404,
				message: "not supported Content-Type\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			middleware := handler.MiddlewareTypeContent(http.HandlerFunc(OKHandler))
			middleware.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}

func TestMiddlewareURLPath(t *testing.T) {
	type want struct {
		code    int
		message string
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "positive test",
			URL:  "/a/b/c/d",
			want: want{
				code:    200,
				message: "",
			},
		},
		{
			name: "negative test [too short path]",
			URL:  "/a/b/c",
			want: want{
				code:    404,
				message: "invalid URL\n",
			},
		},
		{
			name: "negative test [too long path]",
			URL:  "/a/b/c/d/e",
			want: want{
				code:    404,
				message: "invalid URL\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.URL, nil)
			w := httptest.NewRecorder()

			middleware := handler.MiddlewareURLPath(http.HandlerFunc(OKHandler))
			middleware.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}

func TestMiddlewareMetricType(t *testing.T) {
	type want struct {
		code    int
		message string
	}
	tests := []struct {
		name       string
		metricType string
		want       want
	}{
		{
			name:       "positive test [gauge]",
			metricType: "gauge",
			want: want{
				code:    200,
				message: "",
			},
		},
		{
			name:       "positive test [counter]",
			metricType: "counter",
			want: want{
				code:    200,
				message: "",
			},
		},
		{
			name:       "negative test",
			metricType: "other",
			want: want{
				code:    400,
				message: "invalid metric type\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/foo/"+tt.metricType, nil)
			w := httptest.NewRecorder()

			middleware := handler.MiddlewareMetricType(http.HandlerFunc(OKHandler))
			middleware.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}

func TestMiddlewareMetricName(t *testing.T) {
	type want struct {
		code    int
		message string
	}
	tests := []struct {
		name       string
		metricName string
		want       want
	}{
		{
			name:       "positive test",
			metricName: "foo",
			want: want{
				code:    200,
				message: "",
			},
		},
		{
			name:       "negative test",
			metricName: "",
			want: want{
				code:    404,
				message: "invalid metric name\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/foo/bar/"+tt.metricName, nil)
			w := httptest.NewRecorder()

			middleware := handler.MiddlewareMetricName(http.HandlerFunc(OKHandler))
			middleware.ServeHTTP(w, request)

			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.message, string(resBody))
		})
	}
}
