package inmemory

import "errors"

type Item struct {
	Name       string   `json:"name"`
	IntValue   *int64   `json:"int_value,omitempty"`
	FloatValue *float64 `json:"float_value,omitempty"`
}

var (
	errEmptyName      = errors.New("name is empty")
	errEmptyValue     = errors.New("value is empty")
	errUndefinedValue = errors.New("value is undefined")
)

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

func (x Item) validate() error {
	if x.Name == "" {
		return errEmptyName
	}
	if x.hasIntValue() && x.hasFloatValue() {
		return errUndefinedValue
	}
	if !x.hasIntValue() && !x.hasFloatValue() {
		return errEmptyValue
	}
	return nil
}
