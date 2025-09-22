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
	router               *http.ServeMux
	updateCounterService *updateCounterService.Service
	updateGaugeService   *updateGaugeService.Service
}

func New(
	updateCounterService *updateCounterService.Service,
	updateGaugeService *updateGaugeService.Service,
) *API {
	return &API{
		router:               http.NewServeMux(),
		updateCounterService: updateCounterService,
		updateGaugeService:   updateGaugeService,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api API) Route() {
	api.AddRoutes([]Route{
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
	})
}

func (api API) AddRoute(route Route) {
	if route.Method == "" {
		api.router.Handle(route.Path, route.Handler)
	} else {
		api.router.Handle(route.Method+" "+route.Path, route.Handler)
	}
}

func (api API) AddRoutes(routes []Route) {
	for _, r := range routes {
		api.AddRoute(r)
	}
}
