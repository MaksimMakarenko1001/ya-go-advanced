package models

import (
	"time"
)

type OutboxDestination string

const (
	FileOutboxDestination OutboxDestination = "file"
	UrlOutboxDestination  OutboxDestination = "url"
)

type FileEvent struct {
	TS        time.Time `json:"ts"`
	Metrics   []string  `json:"metrics"`
	IPAddress string    `json:"ip_address"`
}

type UrlEvent struct {
	TS        time.Time `json:"ts"`
	Metrics   []string  `json:"metrics"`
	IPAddress string    `json:"ip_address"`
}
