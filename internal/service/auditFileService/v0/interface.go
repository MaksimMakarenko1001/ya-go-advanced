package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/models"
)

type OutboxRepository interface {
	OutboxGetNext(ctx context.Context, destination models.OutboxDestination, segment string, limit int) (resp []entities.Outbox, err error)
	OutboxCommit(ctx context.Context, okIds []entities.OutboxID, failedIds []entities.OutboxID, segment string) (err error)
}

type FileRepository interface {
	FileAppend(ctx context.Context, line []byte) error
}
