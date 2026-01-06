package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	dumpMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/dumpMetricService/v0"
	getFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getFlatService/v0"
	getService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getService/v0"
	hashService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/hashService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateBatchService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateBatchService/v0"
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

	updateFlatService  *updateFlatService.Service
	updateBatchService *updateBatchService.Service
	updateService      *updateService.Service

	getFlatService *getFlatService.Service
	getService     *getService.Service

	listMetricService *listMetricService.Service

	dumpMetricService *dumpMetricService.Service
	hashService       *hashService.Service
}

func New(
	logger logger.HTTPLogger,
	updateFlatService *updateFlatService.Service,
	updateBatchService *updateBatchService.Service,
	updateService *updateService.Service,
	getFlatService *getFlatService.Service,
	getService *getService.Service,
	listMetricService *listMetricService.Service,
	dumpMetricService *dumpMetricService.Service,
	hashService *hashService.Service,
) *API {
	return &API{
		router:             chi.NewRouter(),
		logger:             logger,
		updateFlatService:  updateFlatService,
		updateBatchService: updateBatchService,
		updateService:      updateService,
		getFlatService:     getFlatService,
		getService:         getService,
		listMetricService:  listMetricService,
		dumpMetricService:  dumpMetricService,
		hashService:        hashService,
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

func (api API) HandleIndex(middlewares ...Middleware) {
	api.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/", func(w http.ResponseWriter, rq *http.Request) {
			DoListMetricResponse(api.listMetricService.Do).ServeHTTP(w, rq)
		})
	})
}

func (api API) HandleUpdate(middlewares ...Middleware) {
	api.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, rq *http.Request) {
			DoUpdateFlatResponse(
				api.updateFlatService.Do, chi.URLParam(rq, "type"), chi.URLParam(rq, "name"), chi.URLParam(rq, "value"),
			).ServeHTTP(w, rq)
		})
	})
}

func (api API) HandleUpdateJSON(middlewares ...Middleware) {
	api.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/update/", func(w http.ResponseWriter, rq *http.Request) {
			DoUpdateJSONResponse(api.updateService.Do).ServeHTTP(w, rq)
		})
	})
}

func (api API) HandleUpdateBatchJSON(middlewares ...Middleware) {
	api.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/updates/", func(w http.ResponseWriter, rq *http.Request) {
			DoUpdateBatchJSONResponse(api.updateBatchService.Do).ServeHTTP(w, rq)
		})
	})
}

func (api API) HandleGet(middlewares ...Middleware) {
	api.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/value/{type}/{name}", func(w http.ResponseWriter, rq *http.Request) {
			DoGetFlatResponse(
				api.getFlatService.Do, chi.URLParam(rq, "type"), chi.URLParam(rq, "name"),
			).ServeHTTP(w, rq)
		})
	})
}

func (api API) HandleGetJSON(middlewares ...Middleware) {
	api.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/value/", func(w http.ResponseWriter, rq *http.Request) {
			DoGetJSONResponse(api.getService.Do).ServeHTTP(w, rq)
		})
	})
}

func (api API) WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var resp ResponseInfo
		rw := responseWriter{
			ResponseWriter: w,
			response:       &resp,
		}

		h.ServeHTTP(&rw, r)

		api.logger.LogHTTP(logger.HTTPInfo{
			URI:      r.RequestURI,
			Method:   r.Method,
			Duration: time.Since(start),
			Response: logger.ResponseInfo{
				Size:   resp.Size,
				Status: resp.Status,
				Body:   resp.Body,
			},
		})
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

func (api API) WithHash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hash := r.Header.Get("HashSHA256"); hash != "" {
			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, r.Body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if err := api.hashService.Validate(r.Context(), buf.Bytes(), hash); err != nil {
				WriteError(w, err)
				return
			}

			r.Body = io.NopCloser(buf)
		}

		hw := responseHashWriter{
			ResponseWriter: w,
			body:           bytes.Buffer{},
			hashFunc: func(message []byte) (string, error) {
				return api.hashService.Hash(r.Context(), message)
			},
		}
		h.ServeHTTP(&hw, r)
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
