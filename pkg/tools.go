package pkg

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
