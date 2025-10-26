package config

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/encode"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/storage/inmemory"
	dumpMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/dumpMetricService/v0"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getFlatService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateCounterService/v0"
	updateFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateFlatService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateGaugeService/v0"
	updateService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateService/v0"
)

type DI struct {
	config       *diConfig
	repositories struct {
		metricStorage *inmemory.Repository
		encoder       *encode.JSONEncode
	}
	services struct {
		included struct {
			updateCounterService *updateCounterService.Service
			updateGaugeService   *updateGaugeService.Service
		}
		updateFlatService *updateFlatService.Service
		updateService     *updateService.Service

		getCounterService *getCounterService.Service
		getGaugeService   *getGaugeService.Service
		getFlatService    *getFlatService.Service

		listMetricService *listMetricService.Service

		dumpMetricService *dumpMetricService.Service
	}
	api struct {
		external *handler.API
	}
	infr struct {
		db *db.PGConnect
	}
}

func (di *DI) Init(envPrefix string) {
	di.config = &diConfig{}
	di.config.loadConfig(envPrefix)

	di.initDB(context.Background())
	di.initRepositories()
	di.initServices()
	di.initAPI()
}

func (di *DI) initDB(ctx context.Context) {
	var err error

	initCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	di.infr.db, err = db.New(initCtx, di.config.Database)
	if err != nil {
		log.Printf("init db: %s", err.Error())
	}
}

func (di *DI) initRepositories() {
	di.repositories.encoder = encode.New()
	di.repositories.metricStorage = inmemory.New(di.repositories.encoder)
}

func (di *DI) initServices() {
	di.services.included.updateCounterService = updateCounterService.New(di.repositories.metricStorage)
	di.services.included.updateGaugeService = updateGaugeService.New(di.repositories.metricStorage)

	di.services.updateFlatService = updateFlatService.New(di.services.included.updateCounterService,
		di.services.included.updateGaugeService)
	di.services.updateService = updateService.New(di.services.included.updateCounterService,
		di.services.included.updateGaugeService)

	di.services.getCounterService = getCounterService.New(di.repositories.metricStorage)
	di.services.getGaugeService = getGaugeService.New(di.repositories.metricStorage)
	di.services.getFlatService = getFlatService.New(di.services.getCounterService, di.services.getGaugeService)

	di.services.listMetricService = listMetricService.New(di.repositories.metricStorage)

	di.services.dumpMetricService = dumpMetricService.New(di.config.FileStoragePath, di.repositories.metricStorage)
}

func (di *DI) initAPI() {
	di.api.external = handler.New(
		logger.New(di.config.Logger),
		di.services.updateFlatService,
		di.services.updateService,
		di.services.getCounterService,
		di.services.getGaugeService,
		di.services.getFlatService,
		di.services.listMetricService,
		di.services.dumpMetricService,
	)
	di.api.external.PingHandle(di.infr.db)
	di.api.external.UpdateHandle()
	di.api.external.UpdateJSONHandle(di.config.StoreInterval == 0)
	di.api.external.GetHandle()
}

func (di *DI) Start() error {
	config := di.config.HTTP

	if di.config.Restore {
		if err := di.services.dumpMetricService.ReadDump(); err != nil {
			log.Println(err.Error())
		}
	}

	di.api.external.Route()

	var wg sync.WaitGroup
	errCh := make(chan error)

	if di.config.StoreInterval > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ticker := time.NewTicker(di.config.StoreInterval)
			defer ticker.Stop()

			for {
				select {
				case <-errCh:
					return
				case <-ticker.C:
					if err := di.services.dumpMetricService.WriteDump(); err != nil {
						log.Println(err.Error())
					}
				}
			}
		}()
	}

	err := http.ListenAndServe(config.Address, handler.Conveyor(
		di.api.external,
		di.api.external.WithLogging,
		handler.MiddlewareCompress,
	))

	errCh <- err
	wg.Wait()

	di.infr.db.Close()

	return err
}
