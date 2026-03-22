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
	v0 "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateBatchService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg"
)

func BenchmarkUpdateBatchService(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	now := time.Now()
	counter := models.Metric{
		ID:    "counter",
		MType: pkg.MetricTypeCounter,
		Delta: pkg.ToPtr(int64(10)),
	}
	gauge := models.Metric{
		ID:    "gauge",
		MType: pkg.MetricTypeGauge,
		Value: pkg.ToPtr(99.99),
	}
	metrics := []models.Metric{counter, gauge}

	mockRepo := NewMockMetricRepository(ctrl)
	mockRepo.EXPECT().AddUpdateBatch(
		context.Background(),
		gomock.Eq([]entities.CounterItem{{
			MetricType:  counter.MType,
			MetricName:  counter.ID,
			MetricValue: 10,
			CreatedAt:   now,
			UpdatedAt:   now,
		}}),
		gomock.Eq([]entities.GaugeItem{{
			MetricType:  gauge.MType,
			MetricName:  gauge.ID,
			MetricValue: *gauge.Value,
			CreatedAt:   now,
			UpdatedAt:   now,
		}}),
		gomock.Any(),
		"",
	).AnyTimes().Return(true, nil)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(metrics); err != nil {
		b.Errorf("metrics not ok, %s", err.Error())
	}

	call := func(ctx context.Context, _ time.Time, _ models.Request) (err error) {
		srv := v0.New(mockRepo)
		return srv.Do(ctx, now, models.Request{
			IPAddress: "localhost",
			Metrics:   metrics,
		})
	}
	h := handler.DoUpdateBatchJSONResponse(call)

	for b.Loop() {
		request := httptest.NewRequest(http.MethodPost, "/updates/", &buf)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
	}
}
