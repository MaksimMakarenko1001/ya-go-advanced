package benchmarks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/models"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getCounterService/v0"
	cb "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getCounterService/v0/benchmarks"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getGaugeService/v0"
	gb "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getGaugeService/v0/benchmarks"
	getService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg"
)

func BenchmarkGetService(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	ts := time.Now()
	metric := models.Metric{
		ID:    "counter",
		MType: pkg.MetricTypeCounter,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(metric); err != nil {
		b.Errorf("metric not ok, %s", err.Error())
	}

	mockCounterRepo := cb.NewMockMetricRepository(ctrl)
	mockCounterRepo.EXPECT().GetCounter(context.Background(), gomock.Eq("counter")).AnyTimes().Return(
		&entities.CounterItem{
			MetricType:  pkg.MetricTypeCounter,
			MetricName:  "counter",
			MetricValue: 10,
			CreatedAt:   ts,
			UpdatedAt:   ts,
		}, true, nil)

	mockGaugeRepo := gb.NewMockMetricRepository(ctrl)
	mockGaugeRepo.EXPECT().GetGauge(context.Background(), gomock.Eq("gauge")).AnyTimes().Return(
		&entities.GaugeItem{
			MetricType:  pkg.MetricTypeGauge,
			MetricName:  "gauge",
			MetricValue: 99.99,
			CreatedAt:   ts,
			UpdatedAt:   ts,
		}, true, nil)

	call := func(ctx context.Context, _, _ string) (*models.Metric, error) {
		srv := getService.New(getCounterService.New(mockCounterRepo), getGaugeService.New(mockGaugeRepo))
		return srv.Do(ctx, metric.MType, metric.ID)
	}
	h := handler.DoGetJSONResponse(call)

	for b.Loop() {
		b.StopTimer()
		request := httptest.NewRequest(http.MethodPost, "/value/", &buf)
		w := httptest.NewRecorder()

		b.StartTimer()
		h.ServeHTTP(w, request)
	}
}
