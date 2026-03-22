package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/models"
)

// ExampleDoListMetricResponse demonstrates how to use DoListMetricResponse handler
// to list all metrics in HTML format.
// The handler returns an HTML table with all metrics (both counters and gauges).
func ExampleDoListMetricResponse() {
	// Create the handler with the service
	h := handler.DoListMetricResponse(func(ctx context.Context, template string) (string, error) {
		return "<example/>", nil
	})

	// Use the handler
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Content-Type"))
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// text/html
	// <example/>
}

// ExampleDoGetJSONResponse demonstrates how to use DoGetJSONResponse handler
// to get a metric value in JSON format.
// The handler expects a JSON request body with metric type and ID,
// and returns the metric value in JSON format.
func ExampleDoGetJSONResponse() {
	// Create the handler with the service
	h := handler.DoGetJSONResponse(func(ctx context.Context, metricType, metricName string) (metric *models.Metric, err error) {
		value := 42.5
		return &models.Metric{ID: metricName, MType: metricType, Value: &value}, nil
	})

	// Use the handler with a JSON request body
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id":"test_metric","type":"gauge"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Content-Type"))
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// application/json
	// {"id":"test_metric","type":"gauge","value":42.5}
}

// ExampleDoUpdateJSONResponse demonstrates how to use DoUpdateJSONResponse handler
// to update a metric value in JSON format.
// The handler expects a JSON request body with metric details,
// and returns the updated metric in JSON format.
func ExampleDoUpdateJSONResponse() {
	// Create the handler with the service
	h := handler.DoUpdateJSONResponse(func(ctx context.Context, metric models.Metric) (err error) {
		return nil
	})

	// Use the handler with a JSON request body
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id":"test_metric","type":"gauge","value":42.5}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Content-Type"))
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// application/json
	// {"id":"test_metric","type":"gauge","value":42.5}
}

// ExampleDoUpdateBatchJSONResponse demonstrates how to use DoUpdateBatchJSONResponse handler
// to update multiple metrics in JSON format.
// The handler expects a JSON array of metrics in the request body,
// and returns the updated metrics in JSON format.
func ExampleDoUpdateBatchJSONResponse() {
	// Create the handler with the service
	h := handler.DoUpdateBatchJSONResponse(func(ctx context.Context, ts time.Time, request models.Request) (err error) {
		return nil
	})

	// Use the handler with a JSON array of metrics
	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(`[{"id":"metric1","type":"gauge","value":42.5},{"id":"metric2","type":"counter","delta":10}]`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Content-Type"))
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// application/json
	// [{"id":"metric1","type":"gauge","value":42.5},{"id":"metric2","type":"counter","delta":10}]
}
