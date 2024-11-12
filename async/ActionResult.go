package async

type ActionResult[T any] struct {
	Result T
	Error  error
}

func NewErrorActionResult[T any](err error) ActionResult[T] {
	return NewActionResult(*new(T), err)
}

func NewActionResult[T any](result T, err error) ActionResult[T] {
	return ActionResult[T]{
		Result: result,
		Error:  err,
	}
}
