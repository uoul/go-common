package collections

func Map[T, V any](col []T, fn func(T) V) []V {
	r := make([]V, len(col))
	for i, e := range col {
		r[i] = fn(e)
	}
	return r
}

func Filter[T any](col []T, pred func(T) bool) []T {
	r := []T{}
	for _, e := range col {
		if pred(e) {
			r = append(r, e)
		}
	}
	return r
}
