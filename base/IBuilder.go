package base

type IBuilder[T any] interface {
	Build() T
}
