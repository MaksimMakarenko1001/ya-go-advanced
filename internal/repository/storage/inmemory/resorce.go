package inmemory

type Item struct {
	Name       string   `json:"name"`
	IntValue   *int64   `json:"int_value,omitempty"`
	FloatValue *float64 `json:"float_value,omitempty"`
}

func (x Item) hasIntValue() bool {
	return x.IntValue != nil
}

func (x Item) hasFloatValue() bool {
	return x.FloatValue != nil
}

func (x *Item) add(value int64) {
	value += *x.IntValue
	x.IntValue = &value
}

func (x *Item) update(value float64) {
	x.FloatValue = &value
}
