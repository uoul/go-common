package collections

import "fmt"

func MapSlice[T, V any](col []T, fn func(T) V) []V {
	r := make([]V, len(col))
	for i, e := range col {
		r[i] = fn(e)
	}
	return r
}

func FilterSlice[T any](col []T, pred func(T) bool) []T {
	r := []T{}
	for _, e := range col {
		if pred(e) {
			r = append(r, e)
		}
	}
	return r
}

func ContainsSlice[T any](col []T, pred func(T) bool) bool {
	for _, e := range col {
		if pred(e) {
			return true
		}
	}
	return false
}

func MergeSlice[T1, T2, V any](col1 []T1, col2 []T2, merge func(T1, T2) V) ([]V, error) {
	if len(col1) != len(col2) {
		return nil, fmt.Errorf("collections must be of same length")
	}
	r := make([]V, len(col1))
	for i := 0; i < len(col1); i++ {
		r = append(r, merge(col1[i], col2[i]))
	}
	return r, nil
}

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
