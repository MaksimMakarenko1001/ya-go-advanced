package handler

import (
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateCounterV0Service"
	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateGaugeV0Service"
)

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

type API struct {
	updateCounterService *updateCounterV0Service.Service
	updateGaugeService   *updateGaugeV0Service.Service
}

func New(
	updateCounterService *updateCounterV0Service.Service,
	updateGaugeService *updateGaugeV0Service.Service,
) *API {
	return &API{
		updateCounterService: updateCounterService,
		updateGaugeService:   updateGaugeService,
	}
}

func (api API) Routes() []Route {
	return []Route{
		{
			Method: http.MethodPost,
			Path:   "/update/",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}),
		},
		{
			Method: http.MethodPost,
			Path:   "/update/counter/",
			Handler: Conveyor(DoUpdateCounterResponse(api.updateCounterService.Do),
				MiddlewareTypeContent,
				MiddlewareUrlPath,
				MiddlewareMetricType,
				MiddlewareMetricName,
			),
		},
		{
			Method: http.MethodPost,
			Path:   "/update/gauge/",
			Handler: Conveyor(DoUpdateGaugeResponse(api.updateGaugeService.Do),
				MiddlewareTypeContent,
				MiddlewareUrlPath,
				MiddlewareMetricType,
				MiddlewareMetricName,
			),
		},
	}
}
