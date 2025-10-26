package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	dumpMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/dumpMetricService/v0"
	getFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getFlatService/v0"
	getService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateFlatService/v0"
	updateService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateService/v0"
)

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

type API struct {
	router *chi.Mux
	logger logger.HTTPLogger

	updateFlatService *updateFlatService.Service
	updateService     *updateService.Service

	getFlatService *getFlatService.Service
	getService     *getService.Service

	listMetricService *listMetricService.Service

	dumpMetricService *dumpMetricService.Service
}

func New(
	logger logger.HTTPLogger,
	updateFlatService *updateFlatService.Service,
	updateService *updateService.Service,
	getFlatService *getFlatService.Service,
	getService *getService.Service,
	listMetricService *listMetricService.Service,
	dumpMetricService *dumpMetricService.Service,
) *API {
	return &API{
		router:            chi.NewRouter(),
		logger:            logger,
		updateFlatService: updateFlatService,
		updateService:     updateService,
		getFlatService:    getFlatService,
		getService:        getService,
		listMetricService: listMetricService,
		dumpMetricService: dumpMetricService,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api API) HandlePing(db *db.PGConnect) {
	api.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			WriteError(w, err)
		}

		WriteOK(w)
	})
}

func (api API) HandleIndex() {
	api.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handler := DoListMetricResponse(api.listMetricService.Do)
		handler.ServeHTTP(w, r)
	})
}

func (api API) HandleUpdate() {
	api.router.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handler := DoUpdateFlatResponse(
			api.updateFlatService.Do, chi.URLParam(r, "type"), chi.URLParam(r, "name"), chi.URLParam(r, "value"),
		)

		Conveyor(handler, MiddlewareMetricName).ServeHTTP(w, r)
	})
}

func (api API) HandleUpdateJSON(withSync bool) {
	api.router.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler = DoUpdateJSONResponse(api.updateService.Do)

		if withSync {
			handler = api.WithSync(handler)
		}
		handler.ServeHTTP(w, r)
	})
}

func (api API) HandleGet() {
	api.router.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handler := DoGetFlatResponse(
			api.getFlatService.Do, chi.URLParam(r, "type"), chi.URLParam(r, "name"),
		)

		Conveyor(handler, MiddlewareMetricName).ServeHTTP(w, r)
	})
}

func (api API) HandleGetJSON() {
	api.router.Post("/value/", func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler = DoGetJSONResponse(api.getService.Do)

		handler.ServeHTTP(w, r)
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

		api.logger.LogHTTP(httpInfo)
	})
}

func (api API) WithSync(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		if err := api.dumpMetricService.WriteDump(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
