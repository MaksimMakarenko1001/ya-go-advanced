package logger

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func New(config Config) *ZapLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	logger.Sync()

	lvl, err := zap.ParseAtomicLevel(string(config.Level))
	if err != nil {
		logger.Panic("log level not ok", zap.Error(err))
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		logger.Panic("log config not ok", zap.Error(err))
	}
	logger.Info("zap", zap.String("log level", zl.Level().String()))
	return &ZapLogger{logger: zl}
}

func (zl *ZapLogger) LogHTTP(info HTTPInfo) {
	defer zl.logger.Sync()

	b, err := json.Marshal(info)
	if err != nil {
		zl.logger.Error("log not ok, %w", zap.Error(err))
	}

	infoLabel := "body"
	if info.Response.Status != http.StatusOK {
		infoLabel = "error"
	}

	msg := fmt.Sprint(info.Method, info.URI)

	zl.logger.Info(msg, zap.String(infoLabel, info.Response.Body.String()))
	zl.logger.Debug(msg, zap.ByteString("raw", b))
}
