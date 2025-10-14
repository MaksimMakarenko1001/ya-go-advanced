package handler_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMessage = `Got you`

func testHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(testMessage))
	})

}

func TestMiddlewareTypeContentTextPlain(t *testing.T) {
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
				message: testMessage,
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

			middleware := handler.MiddlewareTypeContentTextPlain(testHandler())
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
				message: testMessage,
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

			middleware := handler.MiddlewareURLPath(testHandler())
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
				message: testMessage,
			},
		},
		{
			name:       "positive test [counter]",
			metricType: "counter",
			want: want{
				code:    200,
				message: testMessage,
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

			middleware := handler.MiddlewareMetricType(testHandler())
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
				message: testMessage,
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

			middleware := handler.MiddlewareMetricName(testHandler())
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

func TestMiddlewareCompress(t *testing.T) {
	middleware := handler.MiddlewareCompress(testHandler())

	request := `What's up?!`

	t.Run("sends gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zw := gzip.NewWriter(buf)

		_, err := zw.Write([]byte(request))
		require.NoError(t, err)

		err = zw.Close()
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/", buf)
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Del("Accept-Encoding")
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, testMessage, string(body))
	})

	t.Run("accept_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(request)
		r := httptest.NewRequest(http.MethodPost, "/", buf)
		r.Header.Set("Accept-Encoding", "gzip")
		r.Header.Del("Content-Encoding")
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, r)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		body, err := io.ReadAll(zr)
		require.NoError(t, err)
		assert.Equal(t, testMessage, string(body))
	})
}
