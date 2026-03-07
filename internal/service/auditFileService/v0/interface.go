package v0

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
)

type OutboxRepository interface {
	OutboxGetNext(ctx context.Context, destination models.OutboxDestination, segment string, limit int) (resp []entities.Outbox, err error)
	OutboxCommit(ctx context.Context, okIds []entities.OutboxID, failedIds []entities.OutboxID, segment string) (err error)
}

type FileRepository interface {
	FileOpen(ctx context.Context) error
	FileAppend(ctx context.Context, line []byte) error
	FileClose(ctx context.Context) error
}
