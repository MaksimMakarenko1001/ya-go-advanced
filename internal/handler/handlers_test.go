package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/handler"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateCounterService/v0"
	updateFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateFlatService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateGaugeService/v0"
	"github.com/stretchr/testify/assert"
)

const html = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <table>
			<tbody>
				<tr>
					<td>counter</td>
					<td>99</td>
				</tr>
				<tr>
					<td>gauge</td>
					<td>99.99</td>
				</tr>
			</tbody>
		</table>
    </body>
</html>`

type MetricRepositoryMock struct {
}

func (m *MetricRepositoryMock) Add(name string, value int64) (ok bool) {
	return name == "ok"
}

func (m *MetricRepositoryMock) Update(name string, value float64) (ok bool) {
	return name == "ok"
}

func (m *MetricRepositoryMock) Get(name string) (any, bool) {
	if name == "ok_counter" {
		return int64(99), true
	}
	if name == "ok_gauge" {
		return float64(99.99), true
	}
	if name == "not_ok_counter" {
		return 99.99, true
	}
	if name == "not_ok_gauge" {
		return "99.99", true
	}
	return nil, false
}

func (m *MetricRepositoryMock) List() []listMetricService.MetricItem {
	return []listMetricService.MetricItem{
		{Name: "gauge", Value: 99.99},
		{Name: "counter", Value: 99},
	}
}

func TestDoListMetricResponse(t *testing.T) {
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name     string
		expected expected
	}{
		{
			name: "positive test",
			expected: expected{
				code: 200,
				body: html,
			},
		},
	}
	handler := handler.DoListMetricResponse(listMetricService.New(&MetricRepositoryMock{}).Do)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, request)

			assert.Equal(t, tt.expected.code, w.Code)
			assert.Equal(t, tt.expected.body, w.Body.String())
		})
	}
}

func TestDoGetCounterResponse(t *testing.T) {
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name       string
		metricName string
		expected   expected
	}{
		{
			name:       "positive test",
			metricName: "ok_counter",
			expected: expected{
				code: 200,
				body: "99",
			},
		},
		{
			name:       "negative test [not found]",
			metricName: "not_found",
			expected: expected{
				code: 404,
				body: "[NOT_FOUND] Not found (`not_found` not found)\n",
			},
		},
		{
			name:       "negative test [exists other type]",
			metricName: "not_ok_counter",
			expected: expected{
				code: 400,
				body: "[BAD_REQUEST] Bad request (`not_ok_counter` type mismatch)\n",
			},
		},
	}
	handler := handler.DoGetCounterResponse(getCounterService.New(&MetricRepositoryMock{}).Do)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/value/counter/"+tt.metricName, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, request)

			assert.Equal(t, tt.expected.code, w.Code)
			if tt.expected.body != "" {
				assert.Equal(t, tt.expected.body, w.Body.String())
			}
		})
	}
}

func TestDoGetGaugeResponse(t *testing.T) {
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name       string
		metricName string
		expected   expected
	}{
		{
			name:       "positive test",
			metricName: "ok_gauge",
			expected: expected{
				code: 200,
				body: "99.99",
			},
		},
		{
			name:       "negative test [not found]",
			metricName: "not_found",
			expected: expected{
				code: 404,
				body: "[NOT_FOUND] Not found (`not_found` not found)\n",
			},
		},
		{
			name:       "negative test [exists other type]",
			metricName: "not_ok_gauge",
			expected: expected{
				code: 400,
				body: "[BAD_REQUEST] Bad request (`not_ok_gauge` type mismatch)\n",
			},
		},
	}
	handler := handler.DoGetGaugeResponse(getGaugeService.New(&MetricRepositoryMock{}).Do)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/value/gauge/"+tt.metricName, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, request)

			assert.Equal(t, tt.expected.code, w.Code)
			if tt.expected.body != "" {
				assert.Equal(t, tt.expected.body, w.Body.String())
			}
		})
	}
}

func TestDoUpdateCounterResponse(t *testing.T) {
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name        string
		metricName  string
		metricValue string
		expected    expected
	}{
		{
			name:        "positive test",
			metricName:  "ok",
			metricValue: "99",
			expected: expected{
				code: 200,
				body: "",
			},
		},
		{
			name:        "negative test [invalide type]",
			metricName:  "ok",
			metricValue: "99.99",
			expected: expected{
				code: 400,
				body: "[BAD_REQUEST] Bad request (invalid metric value)\n",
			},
		},
		{
			name:        "negative test [exists other type]",
			metricName:  "not_ok",
			metricValue: "100",
			expected: expected{
				code: 400,
				body: "",
			},
		},
	}

	service := updateFlatService.New(
		updateCounterService.New(&MetricRepositoryMock{}),
		nil,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/counter/"+tt.metricName+"/"+tt.metricValue, nil)
			w := httptest.NewRecorder()

			handler.DoUpdateFlatResponse(service.Do, "counter", tt.metricName, tt.metricValue).ServeHTTP(w, request)

			assert.Equal(t, tt.expected.code, w.Code)
			if tt.expected.body != "" {
				assert.Equal(t, tt.expected.body, w.Body.String())
			}
		})
	}
}

func TestDoUpdateGaugeResponse(t *testing.T) {
	type expected struct {
		code int
		body string
	}
	tests := []struct {
		name        string
		metricName  string
		metricValue string
		expected    expected
	}{
		{
			name:        "positive test",
			metricName:  "ok",
			metricValue: "99.99",
			expected: expected{
				code: 200,
				body: "",
			},
		},
		{
			name:        "negative test [invalide type]",
			metricName:  "ok",
			metricValue: "99,99",
			expected: expected{
				code: 400,
				body: "[BAD_REQUEST] Bad request (invalid metric value)\n",
			},
		},
		{
			name:        "negative test [exists other type]",
			metricName:  "not_ok",
			metricValue: "100.000",
			expected: expected{
				code: 400,
				body: "",
			},
		},
	}

	service := updateFlatService.New(
		nil,
		updateGaugeService.New(&MetricRepositoryMock{}),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/gauge/"+tt.metricName+"/"+tt.metricValue, nil)
			w := httptest.NewRecorder()

			handler.DoUpdateFlatResponse(service.Do, "gauge", tt.metricName, tt.metricValue).ServeHTTP(w, request)

			assert.Equal(t, tt.expected.code, w.Code)
			if tt.expected.body != "" {
				assert.Equal(t, tt.expected.body, w.Body.String())
			}
		})
	}
}
