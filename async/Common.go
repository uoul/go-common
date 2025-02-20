package async

func Of[T any](task func() (T, error)) chan ActionResult[T] {
	r := make(chan ActionResult[T])
	go func() {
		result, err := task()
		r <- ActionResult[T]{
			Result: result,
			Error:  err,
		}
	}()
	return r
}
