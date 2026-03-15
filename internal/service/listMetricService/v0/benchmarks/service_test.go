package benchmarks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler"
	v0 "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/listMetricService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg"
)

func BenchmarkListMetricService(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	now := time.Now()

	mockRepo := NewMockMetricRepository(ctrl)
	mockRepo.EXPECT().List(context.Background()).AnyTimes().Return(v0.MetricData{
		Counters: []entities.CounterItem{
			{
				MetricType:  pkg.MetricTypeCounter,
				MetricName:  "counter",
				MetricValue: 10,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		Gauges: []entities.GaugeItem{
			{
				MetricType:  pkg.MetricTypeGauge,
				MetricName:  "gauge",
				MetricValue: 99.99,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}, nil)

	call := func(ctx context.Context, _ string) (index string, err error) {
		srv := v0.New(mockRepo)
		return srv.Do(ctx, "")
	}
	h := handler.DoListMetricResponse(call)

	for b.Loop() {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
	}
}
