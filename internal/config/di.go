package config

import (
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/repository/memStorage"
	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateCounterV0Service"
	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateGaugeV0Service"
)

type DI struct {
	config       *diConfig
	repositories struct {
		metricStorage *memStorage.Repository
	}
	services struct {
		updateCounterService *updateCounterV0Service.Service
		updateGaugeService   *updateGaugeV0Service.Service
	}
	api struct {
		external *handler.API
	}
}

func (di *DI) Init() {
	di.loadConfig()
	di.initRepositories()
	di.initServices()
	di.initAPI()
}

func (di *DI) loadConfig() {
	di.config = &diConfig{}
	di.config.loadConfig()
}

func (di *DI) initRepositories() {
	di.repositories.metricStorage = memStorage.New()
}

func (di *DI) initServices() {
	di.services.updateCounterService = updateCounterV0Service.New(di.repositories.metricStorage)
	di.services.updateGaugeService = updateGaugeV0Service.New(di.repositories.metricStorage)
}

func (di *DI) initAPI() {
	di.api.external = handler.New(
		di.services.updateCounterService,
		di.services.updateGaugeService,
	)
}

func (di *DI) Start() error {
	config := di.config.HTTP
	mux := http.NewServeMux()

	for _, route := range di.api.external.Routes() {
		if route.Method == "" {
			mux.Handle(route.Path, route.Handler)
		} else {
			mux.Handle(route.Method+" "+route.Path, route.Handler)
		}
	}

	return http.ListenAndServe(":"+config.Port, mux)
}
