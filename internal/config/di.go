package config

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/logger"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/encode"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/storage/inmemory"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/repository/storage/pg"
	dumpMetricService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/dumpMetricService/v0"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getCounterService/v0"
	getFlatService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getFlatService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getGaugeService/v0"
	getService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/getService/v0"
	hashService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/hashService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/listMetricService/v0"
	updateBatchService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateBatchService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateCounterService/v0"
	updateFlatService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateFlatService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateGaugeService/v0"
	updateService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/updateService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg/backoff"
)

type DI struct {
	config       *diConfig
	logger       *logger.ZapLogger
	repositories struct {
		encoder         *encode.JSONEncode
		inmemoryStorage *inmemory.Repository
		pgStorage       *pg.Repository
	}
	services struct {
		included struct {
			updateCounterService *updateCounterService.Service
			updateGaugeService   *updateGaugeService.Service

			getCounterService *getCounterService.Service
			getGaugeService   *getGaugeService.Service
		}
		updateFlatService  *updateFlatService.Service
		updateBatchService *updateBatchService.Service
		updateService      *updateService.Service

		getFlatService *getFlatService.Service
		getService     *getService.Service

		listMetricService *listMetricService.Service

		dumpMetricService *dumpMetricService.Service
		hashService       *hashService.Service
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

	di.InitLogging()
	di.initDB()
	di.initRepositories()
	di.initServices()
	di.initAPI()
}

func (di *DI) InitLogging() {
	var err error
	di.logger, err = logger.New(di.config.Logger)
	if err != nil {
		log.Println("logger init not ok,", err.Error())
	}
}

func (di *DI) initDB() {
	var err error
	di.infr.db, err = db.New(
		di.config.Database,
		backoff.NewBackoff(di.config.Database.MaxRetries, db.ClassifyPgError),
	)
	if err != nil {
		log.Println("db init not ok,", err.Error())
	}
}

func (di *DI) initRepositories() {
	di.repositories.encoder = encode.New()
	di.repositories.inmemoryStorage = inmemory.New(di.repositories.encoder)
	di.repositories.pgStorage = pg.New(di.infr.db, di.repositories.inmemoryStorage)
}

func (di *DI) initServices() {
	di.services.included.updateCounterService = updateCounterService.New(di.repositories.pgStorage)
	di.services.included.updateGaugeService = updateGaugeService.New(di.repositories.pgStorage)

	di.services.included.getCounterService = getCounterService.New(di.repositories.pgStorage)
	di.services.included.getGaugeService = getGaugeService.New(di.repositories.pgStorage)

	di.services.updateFlatService = updateFlatService.New(di.services.included.updateCounterService,
		di.services.included.updateGaugeService)
	di.services.updateBatchService = updateBatchService.New(di.repositories.pgStorage)
	di.services.updateService = updateService.New(di.services.included.updateCounterService,
		di.services.included.updateGaugeService)

	di.services.getFlatService = getFlatService.New(di.services.included.getCounterService,
		di.services.included.getGaugeService)
	di.services.getService = getService.New(di.services.included.getCounterService,
		di.services.included.getGaugeService)

	di.services.listMetricService = listMetricService.New(di.repositories.pgStorage)

	di.services.dumpMetricService = dumpMetricService.New(di.config.FileStoragePath, di.repositories.inmemoryStorage)
	di.services.hashService = hashService.New(di.config.HashService)
}

func (di *DI) initAPI() {
	di.api.external = handler.New(
		di.logger,
		di.services.updateFlatService,
		di.services.updateBatchService,
		di.services.updateService,
		di.services.getFlatService,
		di.services.getService,
		di.services.listMetricService,
		di.services.dumpMetricService,
		di.services.hashService,
	)
}

func (di *DI) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := di.config.HTTP
	if di.config.Restore {
		if err := di.services.dumpMetricService.ReadDump(); err != nil {
			log.Println(err.Error())
		}
	}

	var optMiddlewares []handler.Middleware
	if di.config.StoreInterval > 0 {
		go di.doDump(ctx)
	} else {
		optMiddlewares = append(optMiddlewares, di.api.external.WithSync)
	}

	di.api.external.HandlePing(di.infr.db)
	di.api.external.HandleIndex()

	di.api.external.HandleGet(handler.MiddlewareMetricName)
	di.api.external.HandleUpdate(handler.MiddlewareMetricName)

	di.api.external.HandleGetJSON()
	di.api.external.HandleUpdateJSON(optMiddlewares...)
	di.api.external.HandleUpdateBatchJSON(optMiddlewares...)

	err := http.ListenAndServe(config.Address, handler.Conveyor(
		di.api.external,
		di.api.external.WithLogging,
		handler.MiddlewareCompress,
		di.api.external.WithHash,
	))

	di.infr.db.Close()

	return err
}

func (di *DI) doDump(ctx context.Context) {
	ticker := time.NewTicker(di.config.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := di.services.dumpMetricService.WriteDump(); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
