package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func New(config Config) *ZapLogger {
	lvl, err := zap.ParseAtomicLevel(string(config.Level))
	if err != nil {
		panic(err)
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &ZapLogger{logger: zl}
}

func (zl *ZapLogger) LogHTTP(message string, info HTTPInfo) {
	defer zl.logger.Sync()

	b, err := json.Marshal(info)
	if err != nil {
		zl.logger.Error("log not ok, %w", zap.Error(err))
	}

	zl.logger.Info(message, zap.ByteString("info", b))
}
