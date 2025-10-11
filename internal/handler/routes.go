package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateCounterService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateGaugeService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

type API struct {
	router               *chi.Mux
	logger               logger.HTTPLogger
	updateCounterService *updateCounterService.Service
	updateGaugeService   *updateGaugeService.Service

	getCounterService *getCounterService.Service
	getGaugeService   *getGaugeService.Service

	listMetricService *listMetricService.Service
}

func New(
	logger logger.HTTPLogger,
	updateCounterService *updateCounterService.Service,
	updateGaugeService *updateGaugeService.Service,
	getCounterService *getCounterService.Service,
	getGaugeService *getGaugeService.Service,
	listMetricService *listMetricService.Service,
) *API {
	return &API{
		router:               chi.NewRouter(),
		logger:               logger,
		updateCounterService: updateCounterService,
		updateGaugeService:   updateGaugeService,
		getCounterService:    getCounterService,
		getGaugeService:      getGaugeService,
		listMetricService:    listMetricService,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api API) Route() {
	api.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler := DoListMetricResponse(api.listMetricService.Do)
		Conveyor(handler, api.WithLogging).ServeHTTP(w, r)
	})

	api.router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler

		switch chi.URLParam(r, "type") {
		case pkg.MetricTypeCounter:
			handler = DoUpdateCounterResponse(api.updateCounterService.Do)
		case pkg.MetricTypeGauge:
			handler = DoUpdateGaugeResponse(api.updateGaugeService.Do)
		}

		if handler == nil {
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}
		Conveyor(handler, api.WithLogging, MiddlewareMetricName).ServeHTTP(w, r)
	})

	api.router.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler

		switch chi.URLParam(r, "type") {
		case pkg.MetricTypeCounter:
			handler = DoGetCounterResponse(api.getCounterService.Do)
		case pkg.MetricTypeGauge:
			handler = DoGetGaugeResponse(api.getGaugeService.Do)
		}

		if handler == nil {
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}
		Conveyor(handler, api.WithLogging, MiddlewareMetricName).ServeHTTP(w, r)
	})
}

func (api API) WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		httpInfo := logger.HTTPInfo{
			URI:    r.RequestURI,
			Method: r.Method,
		}

		rw := responseWriter{
			ResponseWriter: w,
			response:       &httpInfo.Response,
		}
		h.ServeHTTP(&rw, r)

		httpInfo.Duration = time.Since(start)

		api.logger.LogHTTP("http", httpInfo)
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
