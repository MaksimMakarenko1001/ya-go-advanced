package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/logger"
	auditFileService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/auditFileService/v0"
	auditRemoteService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/auditRemoteService/v0"
	decryptService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/decryptService/v0"
	dumpMetricService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/dumpMetricService/v0"
	hashService "github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/hashService/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/worker/sworker"
)

type diConfig struct {
	HTTP               HTTPServerConfig          `envPrefix:"HTTP_"`
	Logger             logger.Config             `envPrefix:"LOGGER_"`
	StoreInterval      time.Duration             `env:"STORE_INTERVAL"`
	FileStoragePath    string                    `env:"FILE_STORAGE_PATH"`
	Restore            bool                      `env:"RESTORE"`
	Database           db.Config                 `envPrefix:"DATABASE"`
	HashService        hashService.Config        `envPrefix:"HASH_SERVICE_"`
	DecryptService     decryptService.Config     `envPrefix:"DECRYPT_SERVICE_"`
	DumpService        dumpMetricService.Config  `envPrefix:"DUMP_SERVICE_"`
	DumpSyncService    dumpMetricService.Config  `envPrefix:"DUMP_SYNC_SERVICE_"`
	AuditFileService   auditFileService.Config   `envPrefix:"AUDIT_FILE_SERVICE_"`
	AuditRemoteService auditRemoteService.Config `envPrefix:"AUDIT_REMOTE_SERVICE_"`
	Worker             struct {
		AuditFile   sworker.Config `envPrefix:"AUDIT_FILE_"`
		AuditRemote sworker.Config `envPrefix:"AUDIT_REMOTE_"`
	} `envPrefix:"WORKER_"`
	AuditFile   string `env:"AUDIT_FILE"`
	AuditRemote string `env:"AUDIT_URL"`
}

func (cfg *diConfig) loadConfig(envPrefix string) {
	cfg.loadDefaults(envPrefix)
	cfg.loadFromArg()
	cfg.loadFromEnv(envPrefix)
	cfg.loadFromEnvToPassTests() // meets tests required

	if cfg.StoreInterval > 0 {
		cfg.DumpSyncService.WriteDumpEnable = false
	}
	if cfg.AuditFile == "" {
		cfg.AuditFileService.AuditEnabled = false
	}
	if cfg.AuditRemote == "" {
		cfg.AuditRemoteService.AuditEnabled = false
	}
	if cfg.DecryptService.CryptoKey == "" {
		cfg.DecryptService.DecryptEnabled = false
	}
}

func (cfg *diConfig) loadDefaults(envPrefix string) {
	if err := env.Parse(cfg, env.Options{Prefix: envPrefix}); err != nil {
		log.Printf("env defaults not ok, %s\n", err.Error())
		return
	}
}

func (cfg *diConfig) loadFromArg() {
	var options struct {
		store int
	}

	flag.StringVar(&cfg.HTTP.Address, "a", `:8080`, "server net address")
	flag.IntVar(&options.store, "i", 0, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", "dump.txt", "dump file path")
	flag.BoolVar(&cfg.Restore, "r", false, "restore dump file on start")
	flag.StringVar(&cfg.Database.DSN, "d", "", "data source name")
	flag.StringVar(&cfg.HashService.Key, "k", "", "hash key")
	flag.StringVar(&cfg.AuditFile, "audit-file", "", "audit file name")
	flag.StringVar(&cfg.AuditRemote, "audit-url", "", "audit full url")
	flag.StringVar(&cfg.DecryptService.CryptoKey, "crypto-key", "", "crypto key path")

	flag.Parse()

	cfg.StoreInterval = time.Second * time.Duration(options.store)
}

func (cfg *diConfig) loadFromEnv(envPrefix string) {
	var config diConfig

	if err := env.Parse(&config, env.Options{Prefix: envPrefix /*, RequiredIfNoDef: true*/}); err != nil {
		log.Printf("env not ok, %s\n", err.Error())
		return
	}

	cfg.Logger = config.Logger
	cfg.Database.MaxRetries = config.Database.MaxRetries
	cfg.HashService.Key = config.HashService.Key

	if address := config.HTTP.Address; address != "" {
		cfg.HTTP.Address = address
	}
	if store := config.StoreInterval; store.String() != "0s" {
		cfg.StoreInterval = store
	}
	if fname := config.FileStoragePath; fname != "" {
		cfg.FileStoragePath = fname
	}
	if restore := config.Restore; restore {
		cfg.Restore = restore
	}
	if auditFile := config.AuditFile; auditFile != "" {
		cfg.AuditFile = auditFile
	}
	if auditRemote := config.AuditRemote; auditRemote != "" {
		cfg.AuditRemote = auditRemote
	}
	if cryptoKey := config.DecryptService.CryptoKey; cryptoKey != "" {
		cfg.DecryptService.CryptoKey = cryptoKey
	}
	if dsn, err := config.Database.ToDSN(); err != nil {
		log.Printf("db config not ok, %s\n", err.Error())
	} else {
		cfg.Database.DSN = dsn
	}
}

func (cfg *diConfig) loadFromEnvToPassTests() {
	if address := os.Getenv("ADDRESS"); address != "" {
		cfg.HTTP.Address = address
	}
	if store, err := strconv.Atoi(os.Getenv("STORE_INTERVAL")); err != nil {
		log.Printf("STORE_INTERVAL env not ok, %s\n", err.Error())
	} else {
		cfg.StoreInterval = time.Second * time.Duration(store)
	}
	if fname := os.Getenv("FILE_STORAGE_PATH"); fname != "" {
		cfg.FileStoragePath = fname
	}
	if restore := os.Getenv("RESTORE"); restore != "" {
		cfg.Restore = restore == "true"
	}
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}
	if key := os.Getenv("KEY"); key != "" {
		cfg.HashService.Key = key
	}
	if auditFile := os.Getenv("AUDIT_FILE"); auditFile != "" {
		cfg.AuditFile = auditFile
	}
	if auditRemote := os.Getenv("AUDIT_URL"); auditRemote != "" {
		cfg.AuditRemote = auditRemote
	}
	if cryptoKey := os.Getenv("CRYPTO_KEY"); cryptoKey != "" {
		cfg.DecryptService.CryptoKey = cryptoKey
	}
}

type HTTPServerConfig struct {
	Address string `env:"ADDRESS"`
}
