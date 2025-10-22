package inmemory

import (
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

func (r *Repository) Add(name string, value int64) bool {
	var zero int64
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = &Item{Name: name, IntValue: &zero}
	}

	item := r.collection[name]
	if !item.hasIntValue() {
		return false
	}

	r.collection[name].add(value)
	return true
}

func (r *Repository) Update(name string, value float64) bool {
	var zero float64
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = &Item{Name: name, FloatValue: &zero}
	}

	item := r.collection[name]
	if !item.hasFloatValue() {
		return false
	}

	r.collection[name].update(value)
	return true
}

func (r *Repository) Get(name string) (any, bool) {
	var value any

	if item, ok := r.collection[name]; ok {
		if item.hasIntValue() {
			value = *item.IntValue
		}
		if item.hasFloatValue() {
			value = *item.FloatValue
		}
	}

	return value, value != nil
}

func (r *Repository) List() []listMetricService.MetricItem {
	res := make([]listMetricService.MetricItem, 0, len(r.collection))
	for name, item := range r.collection {
		var value any
		if item.hasIntValue() {
			value = *item.IntValue
		}
		if item.hasFloatValue() {
			value = *item.FloatValue
		}

		if value != nil {
			res = append(res, listMetricService.MetricItem{
				Name:  name,
				Value: value,
			})
		}
	}
	return res
}

func (r *Repository) Load(b []byte) error {
	var data []Item

	err := r.encoder.Decode(b, &data)
	if err != nil {
		return err
	}

	collection := make(map[string]*Item, len(data))
	for _, item := range data {
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
