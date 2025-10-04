package inmemory

import (
	listMetricService "github.com/MaksimMakarenko1001/ya-go-advanced-sprint-1.git/internal/service/listMetricService/v0"
)

type Repository struct {
	collection map[string]any
}

func New() *Repository {
	return &Repository{
		collection: make(map[string]any),
	}
}

func (r *Repository) Add(name string, value int64) bool {
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = value
		return true
	}

	prev, ok := r.collection[name].(int64)
	if !ok {
		return false
	}

	r.collection[name] = prev + value
	return true
}

func (r *Repository) Update(name string, value float64) bool {
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = value
		return true
	}

	_, ok := r.collection[name].(float64)
	if !ok {
		return false
	}

	r.collection[name] = value
	return true
}

func (r *Repository) Get(name string) (any, bool) {
	if _, ok := r.collection[name]; !ok {
		return nil, false
	}

	return r.collection[name], true
}

func (r *Repository) List() []listMetricService.MetricItem {
	items := make([]listMetricService.MetricItem, 0, len(r.collection))
	for name, value := range r.collection {
		items = append(items, listMetricService.MetricItem{
			Name:  name,
			Value: value,
		})
	}
	return items
}
