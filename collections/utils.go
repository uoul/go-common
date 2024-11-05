package collections

func Map[T, V any](col []T, fn func(T) V) []V {
	r := make([]V, len(col))
	for i, e := range col {
		r[i] = fn(e)
	}
	return r
}
