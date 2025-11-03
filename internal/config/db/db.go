package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGConnect struct {
	db      *sql.DB
	backoff *backoff.Backoff
}

func New(cfg Config, backoff *backoff.Backoff) (conn *PGConnect, err error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return &PGConnect{
		db:      db,
		backoff: backoff,
	}, nil

}

func (pg *PGConnect) Ping(ctx context.Context) error {
	backoff := pg.backoff.WithLinear(time.Second, time.Second*2)
	fn := func(ctx context.Context) error {
		return pg.db.PingContext(ctx)
	}

	return backoff(fn)(ctx)
}

func (pg *PGConnect) QueryWithOneResult(
	ctx context.Context, dst any, query string, args ...any,
) error {
	backoff := pg.backoff.WithLinear(time.Second, time.Second*2)

	fn := func(ctx context.Context) error {
		row := pg.db.QueryRowContext(ctx, query, args...)

		if err := row.Scan(dst); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errors.Join(err, errScanRow)
		}
		return nil
	}

	return backoff(fn)(ctx)
}

func (pg *PGConnect) QueryWithOneResultJSON(
	ctx context.Context, dst any, query string, args ...any,
) error {

	var res []byte

	if err := pg.QueryWithOneResult(ctx, &res, query, args...); err != nil {
		return err
	}

	if len(res) == 0 {
		return errNoData
	}

	if err := json.Unmarshal(res, dst); err != nil {
		return errors.Join(err, errUnmarshal)
	}
	return nil
}

func (pg *PGConnect) Close() error {
	return pg.db.Close()
}
