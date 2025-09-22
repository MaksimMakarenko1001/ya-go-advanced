package handler

import (
	"net/http"

	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateCounterService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateGaugeService/v0"
)

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

type API struct {
	updateCounterService *updateCounterService.Service
	updateGaugeService   *updateGaugeService.Service
}

func New(
	updateCounterService *updateCounterService.Service,
	updateGaugeService *updateGaugeService.Service,
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
				MiddlewareURLPath,
				MiddlewareMetricType,
				MiddlewareMetricName,
			),
		},
		{
			Method: http.MethodPost,
			Path:   "/update/gauge/",
			Handler: Conveyor(DoUpdateGaugeResponse(api.updateGaugeService.Do),
				MiddlewareTypeContent,
				MiddlewareURLPath,
				MiddlewareMetricType,
				MiddlewareMetricName,
			),
		},
	}
}
