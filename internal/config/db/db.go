package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGConnect struct {
	db *sql.DB
}

func New(ctx context.Context, cfg Config) (*PGConnect, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	pg := PGConnect{db: db}
	if err := pg.Ping(ctx); err != nil {
		return nil, err
	}

	return &pg, nil

}

func (pg PGConnect) Ping(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func (pg PGConnect) Close() error {
	return pg.db.Close()
}
