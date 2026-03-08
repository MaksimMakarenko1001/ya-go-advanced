package pkg

import "encoding/json"

func ValuesToList[K comparable, V any](m map[K]V) (l []V) {
	l = make([]V, 0, len(m))
	for _, v := range m {
		l = append(l, v)
	}
	return l
}

// Creates a new pointer to the value.
func ToPtr[T any](value T) *T {
	return &value
}

func JSONMust(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
