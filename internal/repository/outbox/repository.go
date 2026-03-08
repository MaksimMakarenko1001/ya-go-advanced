package outbox

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/models"
)

type Repository struct {
	conn *db.PGConnect
}

func New(conn *db.PGConnect) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) OutboxGetNext(
	ctx context.Context, destination models.OutboxDestination, segment string, limit int,
) (resp []entities.Outbox, err error) {
	err = r.conn.QueryWithOneResultJSON(ctx,
		&resp,
		"select outbox.outbox_get_next(_destination => $1, _segment => $2, _limit => $3)",
		destination, segment, limit,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) OutboxCommit(ctx context.Context, okIds []entities.OutboxID, failedIds []entities.OutboxID, segment string) (err error) {
	if len(okIds) == 0 && len(failedIds) == 0 {
		return nil
	}

	return r.conn.QueryNoResult(ctx,
		"select outbox.outbox_commit(_ok_ids => $1, _failed_ids => $2, _segment => $3)",
		okIds, failedIds, segment,
	)
}
