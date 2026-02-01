package models

import (
	"time"
)

type OutboxDestination string

type FileAuditEvent struct {
	TS        time.Time `json:"ts"`
	Metrics   []string  `json:"metrics"`
	IPAddress string    `json:"ip_address"`
}

type UrlAuditEvent struct {
	TS        time.Time `json:"ts"`
	Metrics   []string  `json:"metrics"`
	IPAddress string    `json:"ip_address"`
}
