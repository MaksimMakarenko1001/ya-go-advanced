package config

import (
	"log"
	"net/http"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/storage/inmemory"
	dumbMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/dumbMetricService/v0"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateCounterService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateGaugeService/v0"
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

		dumbMetricService *dumbMetricService.Service
	}
	api struct {
		external *handler.API
	}
}

func (di *DI) Init(envPrefix string) {
	di.config = &diConfig{}
	di.config.loadConfig(envPrefix)

	di.initRepositories()
	di.initServices()
	di.initAPI()
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

	di.services.dumbMetricService = dumbMetricService.New(di.config.FileStoragePath, di.repositories.metricStorage)
}

func (di *DI) initAPI() {
	di.api.external = handler.New(
		logger.New(di.config.Logger),
		di.services.updateCounterService,
		di.services.updateGaugeService,
		di.services.getCounterService,
		di.services.getGaugeService,
		di.services.listMetricService,
		di.services.dumbMetricService,
	)

}

func (di *DI) Start() error {
	config := di.config.HTTP

	if di.config.Restore {
		di.services.dumbMetricService.ReadDumb()
	}

	withDumb := di.config.StoreInterval == 0
	di.api.external.Route(withDumb)

	errCh := make(chan error)

	if di.config.StoreInterval > 0 {
		go func() {
			ticker := time.NewTicker(di.config.StoreInterval)
			defer ticker.Stop()

			for {
				select {
				case <-errCh:
					return
				case <-ticker.C:
					if err := di.services.dumbMetricService.WriteDumb(); err != nil {
						log.Println(err.Error())
					}
				}
			}
		}()
	}

	err := http.ListenAndServe(config.Address, handler.Conveyor(
		di.api.external,
		handler.MiddlewareCompress,
		di.api.external.WithLogging,
	))

	errCh <- err

	return err
}
