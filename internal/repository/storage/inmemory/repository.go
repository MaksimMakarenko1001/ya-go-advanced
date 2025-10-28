package inmemory

import (
	"context"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/entities"
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced.git/internal/service/listMetricService/v0"
)

type Encoder interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

type Repository struct {
	collection map[string]*Item
	encoder    Encoder
}

func New(encoder Encoder) *Repository {
	return &Repository{
		collection: make(map[string]*Item),
		encoder:    encoder,
	}
}

func (r *Repository) Add(ctx context.Context, name string, value int64) (bool, error) {
	var zero int64
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = &Item{Name: name, IntValue: &zero}
	}

	item := r.collection[name]
	if !item.hasIntValue() {
		return false, nil
	}

	r.collection[name].add(value)
	return true, nil
}

func (r *Repository) Update(ctx context.Context, name string, value float64) (bool, error) {
	var zero float64
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = &Item{Name: name, FloatValue: &zero}
	}

	item := r.collection[name]
	if !item.hasFloatValue() {
		return false, nil
	}

	r.collection[name].update(value)
	return true, nil
}

func (r *Repository) GetCounter(ctx context.Context, name string) (int64, bool, error) {
	item, ok := r.collection[name]
	if !ok || !item.hasIntValue() {
		return 0, false, nil
	}

	return *item.IntValue, true, nil
}

func (r *Repository) GetGauge(ctx context.Context, name string) (float64, bool, error) {
	item, ok := r.collection[name]
	if !ok || !item.hasFloatValue() {
		return 0., false, nil
	}

	return *item.FloatValue, true, nil
}

func (r *Repository) List(ctx context.Context) (listMetricService.MetricData, error) {
	counters := make([]entities.CounterItem, 0, len(r.collection))
	gauges := make([]entities.GaugeItem, 0, len(r.collection))

	for name, item := range r.collection {
		if item.hasIntValue() {
			counters = append(counters, entities.CounterItem{
				MetricName:  name,
				MetricValue: *item.IntValue,
			})
		}

		if item.hasFloatValue() {
			gauges = append(gauges, entities.GaugeItem{
				MetricName:  name,
				MetricValue: *item.FloatValue,
			})
		}
	}
	return listMetricService.MetricData{
		Counters: counters,
		Gauges:   gauges,
	}, nil
}

func (r *Repository) Load(b []byte) error {
	var data []Item

	err := r.encoder.Decode(b, &data)
	if err != nil {
		return err
	}

	collection := make(map[string]*Item, len(data))
	for _, item := range data {
		if err := item.validate(); err != nil {
			return err
		}
		collection[item.Name] = &item
	}

	r.collection = collection
	return nil
}

func (r *Repository) Save() ([]byte, error) {
	data := make([]Item, 0, len(r.collection))

	for _, item := range r.collection {
		data = append(data, *item)
	}

	return r.encoder.Encode(data)
}
