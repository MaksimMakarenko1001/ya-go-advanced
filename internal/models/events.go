package models

import (
	"time"
)

type OutboxDestination string

const (
	FileOutboxDestination   OutboxDestination = "file"
	RemoteOutboxDestination OutboxDestination = "remote"
)

type FileEvent struct {
	TS        time.Time `json:"ts"`
	Metrics   []string  `json:"metrics"`
	IPAddress string    `json:"ip_address"`
}

type RemoteEvent struct {
	TS        time.Time `json:"ts"`
	Metrics   []string  `json:"metrics"`
	IPAddress string    `json:"ip_address"`
}
