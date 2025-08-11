package collections

func MapMap[T comparable, U, V comparable, W any](col map[T]U, fn func(T, U) (V, W)) map[V]W {
	r := map[V]W{}
	for k, v := range col {
		k1, v1 := fn(k, v)
		r[k1] = v1
	}
	return r
}

func FilterMap[K comparable, V any](col map[K]V, pred func(K, V) bool) map[K]V {
	r := map[K]V{}
	for k, v := range col {
		if pred(k, v) {
			r[k] = v
		}
	}
	return r
}

func ContainsMap[K comparable, V any](col map[K]V, pred func(K, V) bool) bool {
	for k, v := range col {
		if pred(k, v) {
			return true
		}
	}
	return false
}
