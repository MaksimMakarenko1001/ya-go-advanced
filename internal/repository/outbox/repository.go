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
		"select outbox.outbox_get_next(_destination => $1, _segment text => $2, _limit => $3)",
		destination, segment, limit,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) OutboxSetFailed(ctx context.Context, ids []entities.OutboxID, segment string) (err error) {
	return r.conn.QueryNoResult(ctx,
		"select outbox.outbox_set_failed(_ids => $1, _segment text => $2)",
		ids, segment,
	)
}

func (r *Repository) OutboxSetCompleted(ctx context.Context, ids []entities.OutboxID, segment string) (err error) {
	return r.conn.QueryNoResult(ctx,
		"select outbox.outbox_set_completed(_ids => $1, _segment text => $2)",
		ids, segment,
	)
}
