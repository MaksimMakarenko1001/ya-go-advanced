package v0

import (
	"context"
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
)

const segment = ""

type Service struct {
	config     Config
	outboxRepo OutboxRepository
	remoteRepo RemoteRepository
}

func New(config Config, outboxRepo OutboxRepository, remoteRepo RemoteRepository) *Service {
	return &Service{
		config:     config,
		outboxRepo: outboxRepo,
		remoteRepo: remoteRepo,
	}
}

func (srv *Service) Do(ctx context.Context) error {
	items, err := srv.outboxRepo.OutboxGetNext(ctx, models.RemoteOutboxDestination, segment, 100)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	completed, failed := make([]entities.OutboxID, 0, len(items)), make([]entities.OutboxID, 0, len(items))

	if !srv.config.AuditEnabled {
		for _, item := range items {
			completed = append(completed, item.ID)
		}
		if err := srv.outboxRepo.OutboxCommit(ctx, completed, failed, segment); err != nil {
			log.Printf("outbox commit failure {dest=%v, err=%v}\n", models.RemoteOutboxDestination, err.Error())
		}
		return nil
	}

	for _, item := range items {
		if err := srv.remoteRepo.RemoteSend(ctx, item.Payload); err != nil {
			log.Printf("outbox failure {dest=%v, id=%v, err=%v}\n", models.RemoteOutboxDestination, item.ID, err.Error())
			failed = append(failed, item.ID)
		} else {
			completed = append(completed, item.ID)
		}
	}

	if err := srv.outboxRepo.OutboxCommit(ctx, completed, failed, segment); err != nil {
		log.Printf("outbox commit failure {dest=%v, err=%v}\n", models.RemoteOutboxDestination, err.Error())
	}

	return nil
}
