package pkg

// SliceFilter selects from slice where the value matches a predicate.
func SliceFilter[S ~[]T, T any](slice S, f func(T) bool) S {
	filtered := make(S, 0, len(slice))

	for _, item := range slice {
		if f(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

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
