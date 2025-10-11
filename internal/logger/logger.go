package logger

import (
	"time"
)

const (
	Trace   LogLevel = "trace"
	Debug   LogLevel = "debug"
	Info    LogLevel = "info"
	Warning LogLevel = "warning"
	Error   LogLevel = "error"
	Fatal   LogLevel = "fatal"
)

type HTTPLogger interface {
	LogHTTP(message string, info HTTPInfo)
}

type Config struct {
	Level LogLevel `env:"LEVEL" envDefault:"info"`
}

type LogLevel = string

type HTTPInfo struct {
	URI      string        `json:"uri"`
	Method   string        `json:"method"`
	Duration time.Duration `json:"duration"`
	Response ResponseInfo  `json:"response"`
}

type ResponseInfo struct {
	Size   int `json:"size"`
	Status int `json:"status"`
}
