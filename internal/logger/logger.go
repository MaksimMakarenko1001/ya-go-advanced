package logger

import (
	"bytes"
	"time"
)

type HTTPLogger interface {
	LogHTTP(info HTTPInfo)
}

type HTTPInfo struct {
	URI      string        `json:"uri"`
	Method   string        `json:"method"`
	Duration time.Duration `json:"duration"`
	Response ResponseInfo  `json:"response"`
}

type ResponseInfo struct {
	Size   int          `json:"size"`
	Status int          `json:"status"`
	Body   bytes.Buffer `json:"-"`
}
