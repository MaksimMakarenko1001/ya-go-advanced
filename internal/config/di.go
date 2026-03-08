package config

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/logger"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/audit/file"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/audit/remote"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/encode"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/outbox"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/storage/inmemory"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/storage/pg"
	auditFileService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/auditFileService/v0"
	auditRemoteService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/auditRemoteService/v0"
	dumpMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/dumpMetricService/v0"
	getCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getCounterService/v0"
	getFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getFlatService/v0"
	getGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getGaugeService/v0"
	getService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/getService/v0"
	hashService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/hashService/v0"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
	updateBatchService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateBatchService/v0"
	updateCounterService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateCounterService/v0"
	updateFlatService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateFlatService/v0"
	updateGaugeService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateGaugeService/v0"
	updateService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/updateService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/worker/sworker"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg/backoff"
)

type DI struct {
	config       *diConfig
	repositories struct {
		encoder         *encode.JSONEncode
		inmemoryStorage *inmemory.Repository
		pgStorage       *pg.Repository
		outbox          *outbox.Repository
		fileAuditor     *file.Repository
		remoteAuditor   *remote.Repository
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

		auditFileService   *auditFileService.Service
		auditRemoteService *auditRemoteService.Service
	}
	workers struct {
		auditFile   *sworker.SimpleWorker
		auditRemote *sworker.SimpleWorker
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

	di.initDB()
	di.initRepositories()
	di.initServices()
	di.initWorkers()
	di.initAPI()
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
	di.repositories.outbox = outbox.New(di.infr.db)
	di.repositories.fileAuditor = file.New(di.config.AuditFile)
	di.repositories.remoteAuditor = remote.New(di.config.AuditRemote)
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

	di.services.auditFileService = auditFileService.New(di.config.AuditFileService, di.repositories.outbox, di.repositories.fileAuditor)
	di.services.auditRemoteService = auditRemoteService.New(di.config.AuditRemoteService, di.repositories.outbox, di.repositories.remoteAuditor)
}

func (di *DI) initWorkers() {
	di.workers.auditFile = sworker.New(
		di.config.Worker.AuditFile,
		"audit_file",
		di.services.auditFileService.Do,
	)
	di.workers.auditRemote = sworker.New(
		di.config.Worker.AuditRemote,
		"audit_remote",
		di.services.auditRemoteService.Do,
	)
}

func (di *DI) initAPI() {
	di.api.external = handler.New(
		logger.New(di.config.Logger),
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

	di.workers.auditFile.Start(ctx)
	di.workers.auditRemote.Start(ctx)

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
