// Module provides common utility functions for data manipulation and conversion.
package pkg

import "encoding/json"

// ValuesToList converts map values to a slice.
func ValuesToList[K comparable, V any](m map[K]V) (l []V) {
	l = make([]V, 0, len(m))
	for _, v := range m {
		l = append(l, v)
	}
	return l
}

// ToPtr creates a new pointer to the value.
func ToPtr[T any](value T) *T {
	return &value
}

// JSONMust marshals a value to JSON and panics on error.
func JSONMust(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
