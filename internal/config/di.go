package config

import (
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/repository/storage/inmemory"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/getCounterService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/getGaugeService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/listMetricService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateCounterService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/updateGaugeService/v0"
)

type DI struct {
	config       *diConfig
	repositories struct {
		metricStorage *inmemory.Repository
	}
	services struct {
		updateCounterService *updateCounterService.Service
		updateGaugeService   *updateGaugeService.Service

		getCounterService *getCounterService.Service
		getGaugeService   *getGaugeService.Service

		listMetricService *listMetricService.Service
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
	di.repositories.metricStorage = inmemory.New()
}

func (di *DI) initServices() {
	di.services.updateCounterService = updateCounterService.New(di.repositories.metricStorage)
	di.services.updateGaugeService = updateGaugeService.New(di.repositories.metricStorage)

	di.services.getCounterService = getCounterService.New(di.repositories.metricStorage)
	di.services.getGaugeService = getGaugeService.New(di.repositories.metricStorage)

	di.services.listMetricService = listMetricService.New(di.repositories.metricStorage)
}

func (di *DI) initAPI() {
	di.api.external = handler.New(
		di.services.updateCounterService,
		di.services.updateGaugeService,
		di.services.getCounterService,
		di.services.getGaugeService,
		di.services.listMetricService,
	)
}

func (di *DI) Start() error {
	config := di.config.HTTP

	di.api.external.Route()
	return http.ListenAndServe(config.Address, di.api.external)
}
