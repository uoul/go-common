package async

type Stream[T any] chan ActionResult[T]

func NewStream[T any]() Stream[T] {
	return make(Stream[T])
}

func NewBufferedStream[T any](bufferSize uint) Stream[T] {
	return make(Stream[T], bufferSize)
}
