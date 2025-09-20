package memStorage

type Repository struct {
	collection map[string]any
}

func New() *Repository {
	return &Repository{collection: make(map[string]any)}
}

func (r *Repository) Add(name string, value int64) (ok bool) {
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

func (r *Repository) Update(name string, value float64) (ok bool) {
	if _, ok := r.collection[name]; !ok {
		r.collection[name] = value
		return true
	}

	_, ok = r.collection[name].(float64)
	if !ok {
		return false
	}

	r.collection[name] = value
	return true
}
