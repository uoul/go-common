package async

func Of[T any](task func(args ...any) (T, error), args ...any) chan ActionResult[T] {
	r := make(chan ActionResult[T])
	go func() {
		result, err := task(args...)
		r <- ActionResult[T]{
			Result: result,
			Error:  err,
		}
	}()
	return r
}
