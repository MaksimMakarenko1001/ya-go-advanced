package pg

import (
	"context"
	"log"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/config/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/repository/storage/inmemory"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
)

type Repository struct {
	conn     *db.PGConnect
	isAlive  bool
	inmemory *inmemory.Repository
}

func New(conn *db.PGConnect, inmemory *inmemory.Repository) *Repository {
	return &Repository{
		conn:     conn,
		isAlive:  checkAlive(conn),
		inmemory: inmemory,
	}
}

func (r *Repository) Add(ctx context.Context, item entities.CounterItem) (ok bool, err error) {
	if !r.isAlive {
		return r.inmemory.Add(ctx, item)
	}

	var updatedNames []string

	err = r.conn.QueryWithOneResultJSON(
		ctx,
		&updatedNames,
		"select metric.counters_upsert(_items => $1)",
		[]entities.CounterItem{item},
	)

	return len(updatedNames) > 0, err
}

func (r *Repository) Update(ctx context.Context, item entities.GaugeItem) (ok bool, err error) {
	if !r.isAlive {
		return r.inmemory.Update(ctx, item)
	}

	var updatedNames []string

	err = r.conn.QueryWithOneResultJSON(
		ctx,
		&updatedNames,
		"select metric.gauges_upsert(_items => $1)",
		[]entities.GaugeItem{item},
	)

	return len(updatedNames) > 0, err
}

func (r *Repository) GetCounter(ctx context.Context, name string) (*entities.CounterItem, bool, error) {
	if !r.isAlive {
		return r.inmemory.GetCounter(ctx, name)
	}

	var items []entities.CounterItem

	err := r.conn.QueryWithOneResultJSON(
		ctx,
		&items,
		"select metric.counters_list_by_metric_names(_metric_names => $1)",
		[]string{name},
	)
	if err != nil {
		return nil, false, err
	}

	if len(items) == 0 {
		return nil, false, nil
	}

	return &items[0], true, nil
}

func (r *Repository) GetGauge(ctx context.Context, name string) (*entities.GaugeItem, bool, error) {
	if !r.isAlive {
		return r.inmemory.GetGauge(ctx, name)
	}

	var items []entities.GaugeItem

	err := r.conn.QueryWithOneResultJSON(
		ctx,
		&items,
		"select metric.gauges_list_by_metric_names(_metric_names => $1)",
		[]string{name},
	)
	if err != nil {
		return nil, false, err
	}

	if len(items) == 0 {
		return nil, false, nil
	}

	return &items[0], true, nil
}

func (r *Repository) List(ctx context.Context) (resp listMetricService.MetricData, err error) {
	if !r.isAlive {
		return r.inmemory.List(ctx)
	}

	err = r.conn.QueryWithOneResultJSON(
		ctx,
		&resp,
		"select metric.list_metrics()",
	)
	return resp, err
}

func checkAlive(conn *db.PGConnect) bool {
	if conn == nil {
		return false
	}

	initCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := conn.Ping(initCtx); err != nil {
		log.Println("db ping not ok,", err.Error())
		return false
	}

	return true
}
