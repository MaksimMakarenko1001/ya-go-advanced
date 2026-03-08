package entities

import (
	"encoding/json"
	"time"
)

type OutboxID = string

type Outbox struct {
	ID          OutboxID        `json:"id"`
	Destination string          `json:"destination"`
	Segment     string          `json:"segment"`
	Payload     json.RawMessage `json:"payload"`
	LockUntil   *time.Time      `json:"lock_until"`
}
